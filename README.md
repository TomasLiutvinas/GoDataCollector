# Go Data Collector
Was supposed to log in and parse HTML, but that turned out to be cancer, due to one input field which looked something like <input type="text" class="wow"> this input held value of an integer that I wasn't able to retrieve using 	"golang.org/x/net/html" and chromedp.
However I found I can call endpoint directly after logging in and it returns JSON Data already, so I guess it is perfect for my purpose, but sad for the learning experience of GO.
