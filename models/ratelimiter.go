package models

import (
	"sync"
	"time"
)

// Visitor struct'ı bir IP adresi için istek sayısını ve zaman damgasını tutar.
type Visitor struct {
	Requests  int
	LastReset time.Time
}

// visitors map'ı tüm ziyaretçilerin bilgilerini tutar.
var visitors = make(map[string]*Visitor)
var mtx sync.Mutex

// getVisitor IP adresine göre ziyaretçiyi döner, eğer yoksa yenisini oluşturur.
func GetVisitor(ip string) *Visitor {
	mtx.Lock()
	defer mtx.Unlock()

	v, exists := visitors[ip]
	if !exists {
		v = &Visitor{
			Requests:  0,
			LastReset: time.Now(),
		}
		visitors[ip] = v
	}
	return v
}

// requestMiddleware her istek geldiğinde çağrılan middleware fonksiyonudur.
