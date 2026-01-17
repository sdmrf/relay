package paths

func (p Paths) Owned() []string {
	return []string{
		p.InstallDir,
		p.DataDir,
		p.BinDir,
		p.CacheDir,
	}
}
