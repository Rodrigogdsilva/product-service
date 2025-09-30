package service

import (
	"context"
	"os"
	"product-service/src/domain"
	"product-service/src/repository"
	"product-service/test_artefacts/seeder"
	"product-service/test_artefacts/stubs"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	pgxUUID "github.com/vgarvardt/pgx-google-uuid/v5"
)

func TestProductService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ProductService Integration Suite")
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

var _ = Describe("ProductService", func() {
	var productService ProductService
	var productRepo repository.ProductRepository
	var testSeeder *seeder.TestSeeder
	var ctx context.Context

	BeforeEach(func() {
		ctx = context.Background()
		productRepo = repository.NewProduct(db)
		productService = NewProductService(productRepo)
		testSeeder = seeder.NewTestSeeder(db)

		_, err := db.Exec(ctx, "TRUNCATE TABLE products RESTART IDENTITY")
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("Creating a new product", func() {
		It("should successfully create a product with a new UUID", func() {
			// Arrange: Dados do novo produto
			name := "Câmara Fantástica"
			description := "Uma câmara com ótima resolução."
			price := 1299.99
			stock := 15

			// Act: Chama o método Create do serviço
			err := productService.Create(ctx, name, description, price, stock)

			// Assert: Verifica se não houve erros
			Expect(err).NotTo(HaveOccurred())

			// Verify: Verifica se o produto foi de fato criado no banco
			products, err := productRepo.ListProducts(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(products).To(HaveLen(1))
			Expect(products[0].Name).To(Equal(name))
			Expect(products[0].ID).NotTo(Equal(uuid.Nil))
		})
	})

	Describe("Updating a product", func() {
		It("should successfully update an existing product", func() {
			// Arrange: Cria um produto para depois o atualizar
			productToUpdate := stubs.NewProductStub().WithName("Produto Antigo").Get()
			Expect(testSeeder.InsertProduct(ctx, productToUpdate)).To(Succeed())

			updatedProduct := &domain.Product{
				ID:          productToUpdate.ID,
				Name:        "Produto Novo e Melhorado",
				Description: "Nova descrição",
				Price:       99.99,
				Stock:       10,
			}

			// Act: Chama o método Update do service
			err := productService.Update(ctx, updatedProduct)

			// Assert: Verifica se não houve erros
			Expect(err).NotTo(HaveOccurred())

			// Verify: Verifica se o produto foi atualizado no banco
			foundProduct, err := productRepo.GetProductByID(ctx, productToUpdate.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(foundProduct.Name).To(Equal("Produto Novo e Melhorado"))
		})
	})
})
