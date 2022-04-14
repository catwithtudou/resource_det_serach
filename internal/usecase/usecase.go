package usecase

import "github.com/google/wire"

var ProvideSet = wire.NewSet(NewUserUsecase, NewDimensionUsecase, NewDocumentUsecase)
