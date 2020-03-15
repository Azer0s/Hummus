package preprocessor

// FileInfo the fileinfo struct for pinpointing errors
type FileInfo struct {
	Start    int
	End      int
	FileName string
}
