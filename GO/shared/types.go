package shared

type ImageData struct {
	Name       string // Nom de l'image
	Data       []byte // Données binaires de l'image
	FilterType int    // Type de filtre à appliquer
}
