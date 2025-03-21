package http

import (
	"ecommerce_clean/db"
	"ecommerce_clean/internals/product/repository"
	"ecommerce_clean/internals/product/usecase"
	"ecommerce_clean/pkgs/middlewares"
	"ecommerce_clean/pkgs/minio"
	"ecommerce_clean/pkgs/redis"
	"ecommerce_clean/pkgs/token"
	"ecommerce_clean/pkgs/validation"
	"github.com/gin-gonic/gin"
)

func Routes(
	r *gin.RouterGroup,
	sqlDB db.IDatabase,
	validator validation.Validation,
	minioClient *minio.MinioClient,
	cache redis.IRedis,
	token token.IMarker,
) {
	productRepository := repository.NewProductRepository(sqlDB)
	productUseCase := usecase.NewProductUseCase(validator, productRepository, minioClient)
	productHandler := NewProductHandler(productUseCase)

	authMiddleware := middlewares.NewAuthMiddleware(token).TokenAuth(cache)

	productRoute := r.Group("/products").Use(authMiddleware)
	{
		productRoute.GET("", productHandler.GetProducts)
		productRoute.GET("/:id", productHandler.GetProduct)
		productRoute.POST("", productHandler.CreateProduct)
		productRoute.PUT("/:id", productHandler.UpdateProduct)
		productRoute.DELETE("/:id", productHandler.DeleteProduct)
	}
}
