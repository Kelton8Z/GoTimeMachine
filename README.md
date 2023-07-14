# GoTimeMachine

This project is inspired by [Rewind](https://www.rewind.ai), which tracks the history of a computer and allows for searching and scrolling back in time.
![image](https://user-images.githubusercontent.com/22567795/224188926-04cec24e-b87b-4d2b-84d5-ed07dfb62b12.png)
![heic_url](https://filesamples.com/samples/image/heic/sample1.heic)

Under the hood, Rewind takes a screenshot every two seconds, OCRs each frame and stores the chunks in user's directory. 
More technical details can be seen in [Tearing down the Rewind app](https://kevinchen.co/blog/rewind-ai-app-teardown/)

This project mainly differentiates from Rewind by managing screen history as a xet repo instead of saving to user directory locally to explore automating with xethub and deduplicating screen recordings.

Implementationwise, this MacOS app is written in Golang through an objective c bridge library, i.e., [macdriver](https://github.com/progrium/macdriver). 
Specifically, goroutines for tracking, pausing and tracing screen captures are added upon the [pomodoro example](https://github.com/progrium/macdriver/tree/main/examples/pomodoro).
Tracking uses ffmpeg to record screen, pausing <font color="red"> cat </font> the newly saved recording <font color="red"> >> </font> the one root recording shown in tracing.

To run the app, 
1. run xetea locally
2. run cas server locally
3. <font color="red">go run main.go </font>
