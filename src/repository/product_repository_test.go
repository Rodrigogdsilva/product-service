package repository

import (
	"context"
	"errors"
	"os"
	"product-service/src/domain"
	"product-service/test_artefacts/seeder"
	"product-service/test_artefacts/stubs"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	pgxUUID "github.com/vgarvardt/pgx-google-uuid/v5"
)

func TestProductRepository(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ProductRepository Integration Suite")
}

var (
	db       *pgxpool.Pool
	pool     *dockertest.Pool
	resource *dockertest.Resource
)

var _ = BeforeSuite(func() {
	var err error
	pool, err = dockertest.NewPool("")
	Expect(err).NotTo(HaveOccurred())

	resource, err = pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres", Tag: "15-alpine",
		Env: []string{"POSTGRES_USER=postgres", "POSTGRES_PASSWORD=secret", "POSTGRES_DB=test_db"},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	Expect(err).NotTo(HaveOccurred())

	err = pool.Retry(func() error {
		dbURL := "postgres://postgres:secret@" + resource.GetHostPort("5432/tcp") + "/test_db?sslmode=disable"

		config, err := pgxpool.ParseConfig(dbURL)
		if err != nil {
			return err
		}

		config.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
			pgxUUID.Register(conn.TypeMap())
			return nil
		}

		db, err = pgxpool.NewWithConfig(context.Background(), config)
		if err != nil {
			return err
		}
		return db.Ping(context.Background())
	})
	Expect(err).NotTo(HaveOccurred())

	migration, err := os.ReadFile("../../database/000001_create_products_table.up.sql")
	Expect(err).NotTo(HaveOccurred())
	_, err = db.Exec(context.Background(), string(migration))
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	Expect(pool.Purge(resource)).To(Succeed())
})

var _ = Describe("ProductRepository", func() {
	var productRepo ProductRepository
	var testSeeder *seeder.TestSeeder
	var ctx context.Context

	BeforeEach(func() {
		ctx = context.Background()
		productRepo = NewProduct(db)
		testSeeder = seeder.NewTestSeeder(db)

		_, err := db.Exec(ctx, "TRUNCATE TABLE products RESTART IDENTITY")
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("Creating a new product", func() {
		Context("with valid data", func() {
			It("should create the product successfully", func() {
				// Arrange: Cria um produto de teste
				product := stubs.NewProductStub().Get()

				// Act: Tenta criar o produto no banco de dados
				err := productRepo.Create(ctx, product)

				// Assert: Verifica se não ocorreu nenhum erro
				Expect(err).NotTo(HaveOccurred())

				// Verify: Busca o produto recém-criado para garantir que ele foi salvo corretamente
				foundProduct, err := productRepo.GetProductByID(ctx, product.ID)
				Expect(err).NotTo(HaveOccurred())
				Expect(foundProduct.ID).To(Equal(product.ID))
				Expect(foundProduct.Name).To(Equal(product.Name))
			})
		})

		Context("when a required parameter is missing", func() {
			It("should return a parameter missing error", func() {
				// Arrange: Cria um produto sem nome
				product := stubs.NewProductStub().WithName("").Get()

				// Act: Tenta criar o produto no banco de dados
				err := productRepo.Create(ctx, product)

				// Assert: Verifica se o erro ocorreu e se é o erro esperado
				Expect(err).To(HaveOccurred())
				Expect(errors.Is(err, domain.ErrParametersMissing)).To(BeTrue())
			})
		})

		Context("with an invalid price", func() {
			It("should return an invalid price error", func() {
				// Arrange: Cria um produto com preço negativo
				product := stubs.NewProductStub().WithPrice(-10.0).Get()

				// Act: Tenta criar o produto no banco de dados
				err := productRepo.Create(ctx, product)

				// Assert: Verifica se o erro ocorreu e se é o erro esperado
				Expect(err).To(HaveOccurred())
				Expect(errors.Is(err, domain.ErrInvalidPrice)).To(BeTrue())
			})
		})
	})

	Describe("Getting a product by ID", func() {
		Context("when the product exists", func() {
			It("should return the correct product", func() {
				// Arrange: Insere um produto de teste no banco
				product := stubs.NewProductStub().Get()
				err := testSeeder.InsertProduct(ctx, product)
				Expect(err).NotTo(HaveOccurred())

				// Act: Busca pelo ID do produto inserido
				foundProduct, err := productRepo.GetProductByID(ctx, product.ID)

				// Assert: Verifica se o produto retornado é o esperado
				Expect(err).NotTo(HaveOccurred())
				Expect(foundProduct).NotTo(BeNil())
				Expect(foundProduct.ID).To(Equal(product.ID))
				Expect(foundProduct.Name).To(Equal(product.Name))
			})
		})

		Context("when the product does not exist", func() {
			It("should return a product not found error", func() {
				// Arrange: Gera um ID não existente
				nonExistentID := stubs.NewProductStub().Get().ID

				// Act: Tenta buscar o produto pelo ID não existente
				_, err := productRepo.GetProductByID(ctx, nonExistentID)

				// Assert: Verifica se o erro ocorreu e se é o erro esperado
				Expect(err).To(HaveOccurred())
				Expect(errors.Is(err, domain.ErrProductNotFound)).To(BeTrue())
			})
		})
	})

	Describe("Listing all products", func() {
		Context("when there are no products", func() {
			It("should return an empty slice", func() {
				// Act: Lista todos os produtos
				products, err := productRepo.ListProducts(ctx)

				// Assert: Verifica se não ocorreu nenhum erro e se a lista está vazia
				Expect(err).NotTo(HaveOccurred())
				Expect(products).NotTo(BeNil())
				Expect(products).To(BeEmpty())
			})
		})

		Context("when there are multiple products", func() {
			It("should return all products", func() {
				// Arrange: Insere 3 produtos
				Expect(testSeeder.InsertProduct(ctx, stubs.NewProductStub().Get())).To(Succeed())
				Expect(testSeeder.InsertProduct(ctx, stubs.NewProductStub().Get())).To(Succeed())
				Expect(testSeeder.InsertProduct(ctx, stubs.NewProductStub().Get())).To(Succeed())

				// Act: Lista todos os produtos
				products, err := productRepo.ListProducts(ctx)

				// Assert: Verifica se não ocorreu nenhum erro e se a lista contém 3 produtos
				Expect(err).NotTo(HaveOccurred())
				Expect(products).To(HaveLen(3))
			})
		})
	})

	Describe("Updating a product", func() {
		It("should update the product details correctly", func() {
			// Arrange: Insere um produto de teste
			originalProduct := stubs.NewProductStub().Get()
			Expect(testSeeder.InsertProduct(ctx, originalProduct)).To(Succeed())

			// Modifica os dados do produto
			updatedProduct := originalProduct
			updatedProduct.Name = "Nome Atualizado"
			updatedProduct.Price = 199.99
			updatedProduct.Stock = 50

			// Act: Tenta atualizar o produto no banco de dados
			err := productRepo.Update(ctx, updatedProduct)

			// Assert: Verifica se não ocorreu nenhum erro
			Expect(err).NotTo(HaveOccurred())

			// Verify: Busca o produto atualizado para garantir que as mudanças foram aplicadas
			foundProduct, err := productRepo.GetProductByID(ctx, originalProduct.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(foundProduct.Name).To(Equal("Nome Atualizado"))
			Expect(foundProduct.Price).To(BeNumerically("==", 199.99))
			Expect(foundProduct.Stock).To(Equal(50))
		})
	})

	Describe("Deleting a product", func() {
		It("should remove the product from the database", func() {
			// Arrange: Insere um produto de teste
			productToDelete := stubs.NewProductStub().Get()
			Expect(testSeeder.InsertProduct(ctx, productToDelete)).To(Succeed())

			// Act: Tenta deletar o produto
			err := productRepo.Delete(ctx, productToDelete.ID)

			// Assert: Verifica se não ocorreu nenhum erro
			Expect(err).NotTo(HaveOccurred())

			// Verify: Tenta buscar o produto deletado, esperando um erro
			_, err = productRepo.GetProductByID(ctx, productToDelete.ID)
			Expect(err).To(HaveOccurred())
			Expect(errors.Is(err, domain.ErrProductNotFound)).To(BeTrue())
		})
	})
})
