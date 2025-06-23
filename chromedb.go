package main

// 	// Il contenuto HTML da renderizzare
// 	htmlContent := `
// 	<!DOCTYPE html>
// <html lang="en" class="bg-white text-black">
// <head>
//   <meta charset="UTF-8" />
//   <meta name="viewport" content="width=device-width, initial-scale=1" />
//   <script src="https://cdn.tailwindcss.com"></script>
//   <style>
//     body {
//       font-family: 'Courier New', Courier, monospace;
//     }
//   </style>
// </head>
// <body class="flex flex-col justify-center items-center h-screen select-none">

//   <main class="text-center space-y-3">
//     <h1 class="text-xl leading-tight tracking-wide">DAJE BYOS</h1>
//     <p class="text-lg">This screen was rendered by BYOS</p>
//     <a href="#" class="underline">Giacomo Marinangeli</a>
//   </main>

//   <footer class="fixed bottom-4 left-4 right-4 max-w-xl mx-auto">
//     <div class="flex items-center justify-center border border-black rounded-md px-3 py-1 text-xs font-mono leading-none whitespace-nowrap bg-white bg-opacity-90">
//       trmnl.gmar.dev
//     </div>
//   </footer>

// </body>
// </html>
//     `
// ctx, cancel := chromedp.NewContext(context.Background())
// defer cancel()

// ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
// defer cancel()

// var buf []byte

// // Usa url.PathEscape per codificare correttamente il contenuto HTML
// dataURL := "data:text/html," + url.PathEscape(htmlContent)

// err := chromedp.Run(ctx,
// 	chromedp.Navigate(dataURL),
// 	chromedp.WaitVisible("body", chromedp.ByQuery),
// 	chromedp.EmulateViewport(800, 480),
// 	chromedp.FullScreenshot(&buf, 90),
// )
// if err != nil {
// 	log.Fatal(err)
// }

// if err := os.WriteFile("generated/output.png", buf, 0644); err != nil {
// 	log.Fatal(err)
// }

// log.Println("Screenshot salvato in output.png")
