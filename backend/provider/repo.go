package provider

import (
	"github.com/google/wire"

	"github.com/bookandmusic/love-girl/internal/repo"
)

var RepoSet = wire.NewSet(
	repo.NewUserRepo,
	repo.NewFileRepo,
	repo.NewAlbumRepo,
	repo.NewAnniversaryRepo,
	repo.NewMomentRepo,
	repo.NewPlaceRepo,
	repo.NewWishRepo,
	repo.NewEntityFileRepo,
	repo.NewSettingRepo,
)
