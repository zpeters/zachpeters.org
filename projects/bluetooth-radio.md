---
name: "Bluetooth Radio"
years: "2017"
summary: "An old clock radio, rebuilt with a Bluetooth amplifier board as a housewarming gift in 2017."
---

In 2017, a good friend was moving from away so we wanted to give her a going-away/ housewarming gift for her new place. She always had a flair for mid-century modern so we decided that something with that style would be the perfect gift

After some searching around, we found the perfect thing. It was an older clock radio, that had just enough _mid-century_ feel. Now to make it special.

<figure>
<label for="mn-fig1" class="margin-toggle">&#8853;</label>
<input type="checkbox" id="mn-fig1" class="margin-toggle" />
<span class="marginnote">Outside view of an old radio</span>
<img src="/assets/img/bluetooth-radio/IMG_20170420_174940.webp" alt="outside view of an old radio">
</figure>

After searching online, i found a few components to start the _transformation_. Most of the "brains" of this project was the [TDA7492P 50W+50W CSR8635 Bluetooth 4.0 Audio Receiver Digital Amplifier Board](https://web.archive.org/web/20181031002148/https://www.ebay.com/itm/TDA7492P-50W-50W-CSR8635-Bluetooth-4-0-Audio-Receiver-Digital-Amplifier-Board-/172141677883?oid=322323937478). I'm not sure if this particular model is available but there are [dozens of similar boards of similar specs](https://www.ebay.com/sch/i.html?_from=R40&_nkw=TDA7492P+50W%2B50W+CSR8635+Bluetooth+4.0+Audio+Receiver+Digital+Amplifier+Board&_sacat=0). This should cost you around $20.00

<figure>
<label for="mn-fig2" class="margin-toggle">&#8853;</label>
<input type="checkbox" id="mn-fig2" class="margin-toggle" />
<span class="marginnote">Another view of the bluetooth module</span>
<img src="/assets/img/bluetooth-radio/IMG_20170420_175129.webp" alt="another view of the bluetooth module">
</figure>

Construction was relatively simple, my skills are sorely lacking and it maybe took an hour to throw everything together. I disconnected the original power plug from the transformer so that was no longer in use. Then i routed the power plug from the ac adaptor that came with the bluetooth unit out the back. Ideally, I wouldn't want the body of the ac adaptor inside the radio housing, but it is probably safer and produces less heat than the original transformer.

<figure>
<label for="mn-fig3" class="margin-toggle">&#8853;</label>
<input type="checkbox" id="mn-fig3" class="margin-toggle" />
<span class="marginnote">The guts and wires of an old radio</span>
<img src="/assets/img/bluetooth-radio/IMG_20170420_175126.webp" alt="the guts and wires of an old radio">
</figure>

The last piece i needed to figure out was essentially converting a stereo signal to mono. It would seem that you can just hook _both_ left and right channels to the "positive" on the speaker and you'd good to go. Unfortunately, this is not the case. You will notice some of the sound comes through but everything sounds distorted. This is due to the fact that the speaker is trying to reproduce two _different_ signals at the same time.

<figure>
<label for="mn-fig4" class="margin-toggle">&#8853;</label>
<input type="checkbox" id="mn-fig4" class="margin-toggle" />
<span class="marginnote">Rear view of a speaker</span>
<img src="/assets/img/bluetooth-radio/IMG_20170420_175140.webp" alt="rear view of a speaker">
</figure>

Thankfully, this is easy to remedy. With a [voltage divider](https://en.wikipedia.org/wiki/Voltage_divider) you can _sum_ the signal, producing a mixed signal. I never ended up writing down the exact components and specifications that I used but you should be able to come up with something by looking up a basic circuit and experimenting with different values.

As you can see the _finished_ build looks more or less the same as the original. The other addition was connecting the top button to the pairing circuit of the bluetooth board. I could never figure out its original function but it worked beautifully for my purposes

<figure>
<label for="mn-fig5" class="margin-toggle">&#8853;</label>
<input type="checkbox" id="mn-fig5" class="margin-toggle" />
<span class="marginnote">Top view of the radio</span>
<img src="/assets/img/bluetooth-radio/IMG_20170420_174945.webp" alt="top view of the radio">
</figure>
