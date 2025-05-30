package handler

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"shoeshop/product-service/internal/model"
	"shoeshop/product-service/internal/service"
	pb "shoeshop/proto"
)

type GRPCHandler struct {
	pb.UnimplementedProductServiceServer
	productService service.ProductService
}

func NewGRPCHandler(productService service.ProductService) *GRPCHandler {
	return &GRPCHandler{
		productService: productService,
	}
}

func (h *GRPCHandler) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.ProductResponse, error) {
	product, err := model.FromProto(req.GetProduct())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid product data: %v", err)
	}

	createdProduct, err := h.productService.CreateProduct(ctx, product)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create product: %v", err)
	}

	return &pb.ProductResponse{
		Product: createdProduct.ToProto(),
	}, nil
}

func (h *GRPCHandler) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.ProductResponse, error) {
	product, err := h.productService.GetProduct(ctx, req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get product: %v", err)
	}

	return &pb.ProductResponse{
		Product: product.ToProto(),
	}, nil
}

func (h *GRPCHandler) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.ProductResponse, error) {
	product, err := model.FromProto(req.GetProduct())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid product data: %v", err)
	}

	updatedProduct, err := h.productService.UpdateProduct(ctx, product)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update product: %v", err)
	}

	return &pb.ProductResponse{
		Product: updatedProduct.ToProto(),
	}, nil
}

func (h *GRPCHandler) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*pb.DeleteProductResponse, error) {
	err := h.productService.DeleteProduct(ctx, req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete product: %v", err)
	}

	return &pb.DeleteProductResponse{
		Success: true,
	}, nil
}

func (h *GRPCHandler) ListProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
	filter := make(map[string]interface{})
	if req.GetCategory() != "" {
		filter["category"] = req.GetCategory()
	}
	if req.GetBrand() != "" {
		filter["brand"] = req.GetBrand()
	}

	products, err := h.productService.ListProducts(ctx, filter)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list products: %v", err)
	}

	pbProducts := make([]*pb.Product, len(products))
	for i, product := range products {
		pbProducts[i] = product.ToProto()
	}

	return &pb.ListProductsResponse{
		Products: pbProducts,
	}, nil
}

func (h *GRPCHandler) SearchProducts(ctx context.Context, req *pb.SearchProductsRequest) (*pb.ListProductsResponse, error) {
	products, err := h.productService.SearchProducts(ctx, req.GetQuery())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to search products: %v", err)
	}

	pbProducts := make([]*pb.Product, len(products))
	for i, product := range products {
		pbProducts[i] = product.ToProto()
	}

	return &pb.ListProductsResponse{
		Products: pbProducts,
	}, nil
} 