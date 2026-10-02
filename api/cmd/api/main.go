package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"gestock/api/internal/categoria"
	"gestock/api/internal/db"
	"gestock/api/internal/deposito"
	"gestock/api/internal/middleware"
	"gestock/api/internal/ordencompra"
	"gestock/api/internal/producto"
	"gestock/api/internal/proveedor"
	"gestock/api/internal/stock"
	"gestock/api/internal/transferencia"
	"gestock/api/internal/usuario"
)

func main() {
	mongoURI := os.Getenv("MONGO_URI")

	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	client, err := db.Conectar(mongoURI)
	if err != nil {
		log.Fatal(err)
	}

	databaseName := os.Getenv("MONGO_DATABASE")

	if databaseName == "" {
		databaseName = "gestock"
	}

	database := client.Database(databaseName)

	jwtSecret := os.Getenv("JWT_SECRET")

	if jwtSecret == "" {
		jwtSecret = "clave-secreta-gestock"
	}

	// Categorías

	categoriaCollection := database.Collection("categorias")

	categoriaRepository := categoria.NewMongoRepository(
		categoriaCollection,
	)

	categoriaService := categoria.NewService(
		categoriaRepository,
	)

	categoriaHandler := categoria.NewHandler(
		categoriaService,
	)

	// Depósitos

	depositoCollection := database.Collection("depositos")

	depositoRepository := deposito.NewMongoRepository(
		depositoCollection,
	)

	depositoService := deposito.NewService(
		depositoRepository,
	)

	depositoHandler := deposito.NewHandler(
		depositoService,
	)

	// Productos

	productoCollection := database.Collection("productos")

	productoRepository := producto.NewMongoRepository(
		productoCollection,
	)

	productoService := producto.NewService(
		productoRepository,
	)

	productoHandler := producto.NewHandler(
		productoService,
	)

	// Stock

	stockRepository := stock.NewMongoRepository(
		database,
	)

	stockService := stock.NewService(
		stockRepository,
	)

	stockHandler := stock.NewHandler(
		stockService,
	)

	// Proveedores

	proveedorCollection := database.Collection("proveedores")

	proveedorRepository := proveedor.NewMongoRepository(
		proveedorCollection,
	)

	proveedorService := proveedor.NewService(
		proveedorRepository,
	)

	proveedorHandler := proveedor.NewHandler(
		proveedorService,
	)

	// Ordenes de compra

	ordenCompraCollection := database.Collection("ordenes_compra")

	ordenCompraRepository := ordencompra.NewMongoRepository(
		ordenCompraCollection,
	)

	ordenCompraService := ordencompra.NewService(
		ordenCompraRepository,
		productoRepository,
		stockService,
	)

	ordenCompraHandler := ordencompra.NewHandler(
		ordenCompraService,
	)

	// Transferencias

	transferenciaRepository := transferencia.NewMongoRepository(
		database,
	)

	transferenciaService := transferencia.NewService(
		transferenciaRepository,
		stockService,
	)

	transferenciaHandler := transferencia.NewHandler(
		transferenciaService,
	)

	// Usuarios

	usuarioCollection := database.Collection("usuarios")

	usuarioRepository := usuario.NewMongoRepository(
		usuarioCollection,
	)

	usuarioService := usuario.NewService(
		usuarioRepository,
		jwtSecret,
	)

	usuarioHandler := usuario.NewHandler(
		usuarioService,
	)

	// Router

	router := gin.Default()

	categoria.RegisterRoutes(
		router,
		categoriaHandler,
	)

	deposito.RegisterRoutes(
		router,
		depositoHandler,
	)

	producto.RegisterRoutes(
		router,
		productoHandler,
	)

	usuario.RegisterRoutes(
		router,
		usuarioHandler,
	)

	// Rutas protegidas

	protected := router.Group("/api")

	protected.Use(
		middleware.Auth(jwtSecret),
	)

	protected.GET("/perfil", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"user_id":     c.GetString("user_id"),
			"rol":         c.GetString("rol"),
			"deposito_id": c.GetString("deposito_id"),
		})
	})

	log.Println("API escuchando en el puerto 8080")

	stock.RegisterRoutes(
		protected,
		stockHandler,
	)

	proveedor.RegisterRoutes(
		protected,
		proveedorHandler,
	)

	ordencompra.RegisterRoutes(
		protected,
		ordenCompraHandler,
	)

	transferencia.RegisterRoutes(
		protected,
		transferenciaHandler,
	)

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
