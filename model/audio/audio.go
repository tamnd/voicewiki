package audio

type Audio struct {
	Server string `gorethink:"server" json:"-"`
	Path   string `gorethink:"path" json:"-"`
	URL    string `gorethink:"-" json:"url"`
}
