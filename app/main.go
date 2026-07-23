package main

//embmodel "neuro/lim/EmbModel"

func main() {
	shundown := RunApp()
	//defer logger.Log.Shutdown()
	defer shundown()
}
