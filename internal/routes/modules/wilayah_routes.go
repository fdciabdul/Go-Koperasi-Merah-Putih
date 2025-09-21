package modules

import (
	"github.com/gin-gonic/gin"
	"koperasi-merah-putih/internal/handlers"
)

type WilayahRoutes struct {
	wilayahHandler *handlers.WilayahHandler
}

func NewWilayahRoutes(wilayahHandler *handlers.WilayahHandler) *WilayahRoutes {
	return &WilayahRoutes{
		wilayahHandler: wilayahHandler,
	}
}

func (r *WilayahRoutes) SetupRoutes(router *gin.RouterGroup) {
	wilayah := router.Group("/wilayah")
	{
		wilayah.GET("/provinsi", r.wilayahHandler.GetProvinsiList)
		wilayah.GET("/provinsi/:provinsi_id/kabupaten", r.wilayahHandler.GetKabupatenList)
		wilayah.GET("/kabupaten/:kabupaten_id/kecamatan", r.wilayahHandler.GetKecamatanList)
		wilayah.GET("/kecamatan/:kecamatan_id/kelurahan", r.wilayahHandler.GetKelurahanList)
	}
}