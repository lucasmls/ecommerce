package repositories

import (
	"context"
	"reflect"
	"testing"

	"github.com/lucasmls/ecommerce/services/products/domain"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func TestMustNewInMemoryProductsRepository(t *testing.T) {
	loggerM, _ := zap.NewDevelopment()
	tracerM := trace.NewNoopTracerProvider().Tracer("")

	t.Run("Failure tests", func(t *testing.T) {
		tt := []struct {
			name  string
			input ProductsRepositoryInput
			want  error
		}{
			{
				name: "Should not be able to instantiate ProductsRepository with no storage capacity",
				input: ProductsRepositoryInput{
					Logger: loggerM,
					Tracer: tracerM,
					Size:   0,
				},
				want: ErrInvalidStorageSize,
			},
		}

		for _, tc := range tt {
			t.Run(tc.name, func(t *testing.T) {
				defer func() {
					err := recover()

					if !reflect.DeepEqual(err, tc.want) {
						t.Errorf("MustNewInMemoryProductsRepository() got = %v, want %v", err, tc.want)
					}
				}()

				MustNewInMemoryProductsRepository(tc.input)
			})
		}
	})

	t.Run("Successful tests", func(t *testing.T) {
		tt := []struct {
			name  string
			input ProductsRepositoryInput
			want  InMemoryProductsRepository
		}{
			{
				name: "Should construct ProductsRepository with correct storage size",
				input: ProductsRepositoryInput{
					Logger: loggerM,
					Tracer: tracerM,
					Size:   10,
				},
				want: InMemoryProductsRepository{
					in: ProductsRepositoryInput{
						Logger: loggerM,
						Tracer: tracerM,
						Size:   10,
					},
					storage: map[string]domain.Product{
						"10": {
							ID:          "10",
							Name:        "Macbook Air M1",
							Description: "Fast!",
							Price:       6900,
						},
					},
				},
			},
		}

		for _, tc := range tt {
			t.Run(tc.name, func(t *testing.T) {
				got := MustNewInMemoryProductsRepository(tc.input)
				if !reflect.DeepEqual(got, tc.want) {
					t.Errorf("NewInMemoryProductsRepository() got = %v, want %v", got, tc.want)
				}
			})
		}
	})
}

func TestNewInMemoryProductsRepository(t *testing.T) {
	loggerM, _ := zap.NewDevelopment()
	tracerM := trace.NewNoopTracerProvider().Tracer("")

	t.Run("Failure tests", func(t *testing.T) {
		tt := []struct {
			name  string
			input ProductsRepositoryInput
			want  error
		}{
			{
				name: "Should not be able to instantiate ProductsRepository with no storage capacity",
				input: ProductsRepositoryInput{
					Logger: loggerM,
					Tracer: tracerM,
					Size:   0,
				},
				want: ErrInvalidStorageSize,
			},
		}

		for _, tc := range tt {
			t.Run(tc.name, func(t *testing.T) {
				_, err := NewInMemoryProductsRepository(tc.input)
				if !reflect.DeepEqual(err, tc.want) {
					t.Errorf("NewInMemoryProductsRepository() got = %v, want %v", err, tc.want)
				}
			})
		}
	})

	t.Run("Successful tests", func(t *testing.T) {
		tt := []struct {
			name  string
			input ProductsRepositoryInput
			want  InMemoryProductsRepository
		}{
			{
				name: "Should construct ProductsRepository with correct storage size",
				input: ProductsRepositoryInput{
					Logger: loggerM,
					Tracer: tracerM,
					Size:   10,
				},
				want: InMemoryProductsRepository{
					in: ProductsRepositoryInput{
						Logger: loggerM,
						Tracer: tracerM,
						Size:   10,
					},
					storage: map[string]domain.Product{
						"10": {
							ID:          "10",
							Name:        "Macbook Air M1",
							Description: "Fast!",
							Price:       6900,
						},
					},
				},
			},
		}

		for _, tc := range tt {
			t.Run(tc.name, func(t *testing.T) {
				got, err := NewInMemoryProductsRepository(tc.input)
				if err != nil {
					t.Errorf("NewInMemoryProductsRepository should not have failed. Received: %v", err)
				}

				if !reflect.DeepEqual(got, tc.want) {
					t.Errorf("NewInMemoryProductsRepository() got = %v, want %v", got, tc.want)
				}
			})
		}
	})
}

func TestInMemoryProductsRepository_Create(t *testing.T) {}

func TestInMemoryProductsRepository_Delete(t *testing.T) {
	loggerM, _ := zap.NewDevelopment()
	tracerM := trace.NewNoopTracerProvider().Tracer("")

	type fields struct {
		in      ProductsRepositoryInput
		storage map[string]domain.Product
	}
	type args struct {
		ctx context.Context
		id  string
	}

	t.Run("Failure tests", func(t *testing.T) {
		tests := []struct {
			name   string
			fields fields
			args   args
			want   error
		}{
			{
				name: "Should return a not found error in case the specified product isn't stored",
				fields: fields{
					in: ProductsRepositoryInput{
						Logger: loggerM,
						Tracer: tracerM,
						Size:   1,
					},
					storage: map[string]domain.Product{},
				},
				args: args{
					ctx: context.Background(),
					id:  "1",
				},
				want: ErrProductNotFound,
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				r := InMemoryProductsRepository{
					in:      tt.fields.in,
					storage: tt.fields.storage,
				}
				err := r.Delete(tt.args.ctx, tt.args.id)
				if !reflect.DeepEqual(err, tt.want) {
					t.Errorf("Delete() got = %v, want %v", err, tt.want)
				}
			})
		}
	})

	t.Run("Successful tests", func(t *testing.T) {
		tests := []struct {
			name   string
			fields fields
			args   args
			want   error
		}{
			{
				name: "Should remove the specified Product",
				fields: fields{
					in: ProductsRepositoryInput{
						Logger: loggerM,
						Tracer: tracerM,
						Size:   2,
					},
					storage: map[string]domain.Product{
						"1": {
							ID:          "1",
							Name:        "Macbook Pro, M1 Max",
							Description: "Incredible!",
							Price:       27600,
						},
					},
				},
				args: args{
					ctx: context.Background(),
					id:  "1",
				},
				want: nil,
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				r := InMemoryProductsRepository{
					in:      tt.fields.in,
					storage: tt.fields.storage,
				}

				err := r.Delete(tt.args.ctx, tt.args.id)
				if !reflect.DeepEqual(err, tt.want) {
					t.Errorf("Delete() got = %v, want %v", err, tt.want)
				}
			})
		}
	})
}

func TestInMemoryProductsRepository_List(t *testing.T) {
	loggerM, _ := zap.NewDevelopment()
	tracerM := trace.NewNoopTracerProvider().Tracer("")

	type fields struct {
		in      ProductsRepositoryInput
		storage map[string]domain.Product
	}
	type args struct {
		ctx    context.Context
		filter domain.ListProductsFilter
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []domain.Product
	}{
		{
			name: "Should list all products that were stored",
			fields: fields{
				in: ProductsRepositoryInput{
					Logger: loggerM,
					Tracer: tracerM,
					Size:   10,
				},
				storage: map[string]domain.Product{
					"1": {
						ID:          "1",
						Name:        "Macbook Pro M1",
						Description: "Fast",
						Price:       8900,
					},
				},
			},
			args: args{
				ctx: context.Background(),
				filter: domain.ListProductsFilter{
					IDs: []string{},
				},
			},
			want: []domain.Product{
				{
					ID:          "1",
					Name:        "Macbook Pro M1",
					Description: "Fast",
					Price:       8900,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := InMemoryProductsRepository{
				in:      tt.fields.in,
				storage: tt.fields.storage,
			}
			got, err := r.List(tt.args.ctx, tt.args.filter)
			if err != nil {
				t.Errorf("List() should not have failed. Received: %v", err)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("List() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInMemoryProductsRepository_Update(t *testing.T) {
	loggerM, _ := zap.NewDevelopment()
	tracerM := trace.NewNoopTracerProvider().Tracer("")

	type fields struct {
		in      ProductsRepositoryInput
		storage map[string]domain.Product
	}

	type args struct {
		ctx     context.Context
		product domain.Product
	}

	t.Run("Failure test", func(t *testing.T) {
		tests := []struct {
			name   string
			fields fields
			args   args
			want   error
		}{
			{
				name: "Should return a not found error in case the specified product isn't stored",
				fields: fields{
					in: ProductsRepositoryInput{
						Logger: loggerM,
						Tracer: tracerM,
						Size:   1,
					},
					storage: map[string]domain.Product{},
				},
				args: args{
					ctx: context.Background(),
					product: domain.Product{
						ID:          "1",
						Name:        "Macbook Pro M1",
						Description: "Fast",
						Price:       8900,
					},
				},
				want: ErrProductNotFound,
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				r := InMemoryProductsRepository{
					in:      tt.fields.in,
					storage: tt.fields.storage,
				}
				_, err := r.Update(tt.args.ctx, tt.args.product)
				if !reflect.DeepEqual(err, tt.want) {
					t.Errorf("Delete() got = %v, want %v", err, tt.want)
				}
			})
		}
	})

	t.Run("Successful test", func(t *testing.T) {
		tests := []struct {
			name   string
			fields fields
			args   args
			want   domain.Product
		}{
			{
				name: "Should update the specified Product correctly",
				fields: fields{
					in: ProductsRepositoryInput{
						Logger: loggerM,
						Tracer: tracerM,
						Size:   2,
					},
					storage: map[string]domain.Product{
						"1": {
							ID:          "1",
							Name:        "Wrong macbook name",
							Description: "Slow",
							Price:       1000,
						},
					},
				},
				args: args{
					ctx: context.Background(),
					product: domain.Product{
						ID:          "1",
						Name:        "Macbook Pro M1",
						Description: "Fast",
						Price:       8900,
					},
				},
				want: domain.Product{
					ID:          "1",
					Name:        "Macbook Pro M1",
					Description: "Fast",
					Price:       8900,
				},
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				r := InMemoryProductsRepository{
					in:      tt.fields.in,
					storage: tt.fields.storage,
				}
				got, err := r.Update(tt.args.ctx, tt.args.product)
				if err != nil {
					t.Errorf("Update() should not have failed. Received: %v", err)
				}

				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Update() got = %v, want %v", got, tt.want)
				}
			})
		}
	})

}
