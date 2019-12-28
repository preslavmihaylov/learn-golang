# Exercise #18: Image Transformer

[![exercise status: released](https://img.shields.io/badge/exercise%20status-released-green.svg?style=for-the-badge)](https://gophercises.com/exercises/transform)

## Exercise details

Create a web server where the user uploads an image, and then is guided through a selection process using various options on the [primitive](https://github.com/fogleman/primitive) CLI.

For instance, a user might upload an image and then be present with outputs using several of the different modes:

```
mode: 0=combo, 1=triangle, 2=rect, 3=ellipse, 4=circle, 5=rotatedrect, 6=beziers, 7=rotatedellipse, 8=polygon
```

A user could then selec the image with the mode they prefer, at which time the web server would then assume the user prefers that mode and give them new options with that mode and another variable, such as the number of shapes (the `n` flag). At this time the server would output maybe 6 samples each using a different `n` value on the original image with the selected mode and the user could choose the they prefer.

While primitive doesn't have a ton of options, there are enough to at least create 2-3 steps like described above and then once all the settings are selected you could produce perhaps 4 images using those settings and let the user choose their favorite amongst those and download it.

### Notes and suggestions

**Use the `os/exec` package**

You should interact with [primitive](https://github.com/fogleman/primitive) via the command line so you can get some practice using the [os/exec](https://golang.org/pkg/os/exec/) package, temporary files, and more. You don't have to do this (*obviously!*) but the reason I chose this exercise was to give you a chance to use these packages and tools.

**Focus on getting things working, then improve**

Your first version will probably produce a lot of images that never get deleted or cleaned up. Don't worry about that at first. Just get your code working.

Once it works try to think of ways that you could manage cleaning up images. For instance, maybe you could keep an in-memory list of every image you have created and a time when it should be deleted (perhaps 5 minutes after it is created?). Or maybe you want to practice using a database so you could store this info there.

I will probably use an in-memory solution and set my timer based on whether or not an image was selected. For instance, I could start all images out with a 5m deletion timer, but if the user selects the `circle` shape image I could update its timer to last an extra 15 minutes since this is the one the user preferred, and I could continue doing this through each step of my app - add time to the selected image, and make sure any unselected images are set to be deleted in 5 minutes.

## Bonus

As a bonus exercise, add support for more image transformation options. For instance, the [legoizer](https://github.com/esimov/legoizer) transformation is a pretty neat one that make an image look like it was printed on lego blocks. You could also write a custom service that will turn an image into pure black and white (not greyscale) so that it could be printed on something like a [Zebra label printer](https://www.zebra.com/us/en/products/printers.html).

