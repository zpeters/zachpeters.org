---
name: "Morse Code Glove"
years: "2016"
date: "2026-10-26T19:04:19-0500"
summary: "A glove that turns finger taps into Morse letters."
---

This is one of my experiments into _eyes-free_[^1] interfaces. The original idea was for a car computer. That project never took off. The arrival of smart phones with _pretty good_ voice assistants removed most of the need.

Nevertheless, I decided to see this one through. Please excuse my primitive construction skills!

<figure>
<img src="/assets/img/morse-code-glove/20160724-IMG_20160724_085317-thumb.webp">
</figure>

The core of this project came together pretty quickly. Construction as you can see is very rough, but it was enough to get the basic concept going. Within a few minutes of its creation I could easily switch tracks and adjust the volume on my media player with little mental interruption of the task I was working on.

<figure>
<img src="/assets/img/morse-code-glove/20160724-IMG_20160724_085127-thumb.webp">
</figure>

Future experimentation would be around the glove. As you can imagine it’s very cumbersome!

### Hardware Interface

This just decodes finger touches (dit and dah) into letters just as if you had typed them into a keyboard. This makes it very simple to interface with other software. In my experimentation I created a simple script menu[^2] to run different commands to control a music player.

### Construction

<figure class="fullwidth">
<table class="parts">
<thead>
<tr>
<th>Part</th>
<th>Description</th>
<th class="num">Cost</th>
<th>URL</th>
</tr>
</thead>
<tbody>
<tr>
<td>Elastic band</td>
<td>Just something simple to mount the board on</td>
<td class="num">$10.00</td>
<td><a href="http://amzn.com/B00SN6987Y">http://amzn.com/B00SN6987Y</a></td>
</tr>
<tr>
<td>Glove</td>
<td>The cheapest conductive glove I could find at Walmart</td>
<td class="num">$3.95</td>
<td><a href="https://goo.gl/pWZz7M">http://goo.gl/pWZz7M</a></td>
</tr>
<tr>
<td>Conductive Thread</td>
<td>You'll need a decently sturdy needle for this.</td>
<td class="num">$8.00</td>
<td><a href="http://amzn.com/B00W3D3X9W">http://amzn.com/B00W3D3X9W</a></td>
</tr>
<tr>
<td>Jumper Wires</td>
<td>Easy way to connect things for experimenting</td>
<td class="num">$5.00</td>
<td><a href="http://amzn.com/B0002KRABU">http://amzn.com/B0002KRABU</a></td>
</tr>
<tr>
<td>Headers</td>
<td>Makes it easy to remove the teensy when you are done tinkering</td>
<td class="num">$5.00</td>
<td><a href="https://goo.gl/JszxBU">http://goo.gl/JszxBU</a></td>
</tr>
<tr>
<td>Perfboard</td>
<td>Really easy way to mount your projects</td>
<td class="num">$8.50</td>
<td><a href="https://www.adafruit.com/product/589">https://www.adafruit.com/product/589</a></td>
</tr>
<tr>
<td>Teensy</td>
<td>The best "arduino like" out there</td>
<td class="num">$20.00</td>
<td><a href="https://www.pjrc.com/store/teensy32.html">https://www.pjrc.com/store/teensy32.html</a></td>
</tr>
</tbody>
<tfoot>
<tr>
<td><em>Total</em></td>
<td></td>
<td class="num"><em>$60.45</em></td>
<td></td>
</tr>
</tfoot>
</table>
<figcaption>
Table 1: Parts for the glove and their cost. Prices as of 9/27/2016.
</figcaption>
</figure>

[^1]: [Oakley, I., Park, JS. (2007). Designing Eyes-Free Interaction. In: Oakley, I., Brewster, S. (eds) Haptic and Audio Interaction Design. HAID 2007. Lecture Notes in Computer Science, vol 4813. Springer, Berlin, Heidelberg. https://doi.org/10.1007/978-3-540-76702-2_13](http://link.springer.com/10.1007/978-3-540-76702-2_13)
[^2]: [https://gist.github.com/zpeters/59a5f6345897309da9eef729780c09f3](https://gist.github.com/zpeters/59a5f6345897309da9eef729780c09f3)
