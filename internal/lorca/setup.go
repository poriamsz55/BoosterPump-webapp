package lorca

func Setup(port string) {
	// folder := tmpFolder() Solve issue on older windows but chrome report error
	// https://github.com/GoogleChrome/chrome-launcher/blob/master/docs/chrome-flags-for-tools.md
	args := []string{
		"--allow-running-insecure-content",
		"--disable-notifications",
		"--allow-insecure-localhost",
	}
	ui, err := New("http://127.0.0.1:"+port+"/", "", 1600, 900, args...) // folder
	if err != nil {
		PromptDownload()
		return
	}
	defer ui.Close()
	_ = ui.SetBounds(Bounds{WindowState: WindowStateMaximized})
	<-ui.Done()
}
