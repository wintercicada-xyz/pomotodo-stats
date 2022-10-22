# Pomotodo Stats

Generate a calendar heatmap according to your [Pomotodo](https://pomotodo.com) pomo history calendar.



## Usage
> :warning: Pomo history calendar is a Pomotodo Pro feature, if you use the free version of Pomotodo, we can't get your pomo history and generate the calendar heatmap.

At first, you should get your pomo history link, you can do that by going to [Pomotodo Web](https://pomotodo.com/app/), opening Preferences, selecting Subscribe, and then saving the link of Pomo history.
![](/images/pomotodo.png)

After that, you can use the link of Pomo history to get the calendar heatmap by using the link ~~https://api.wintercicada.xyz/pomotodo-stats/ics=<the-pomo-history-link>~~ (unmaintained). And then you can get an SVG file.
![](/images/pomotodoHeatmap.svg)



## Privacy
Notice that your pomo history is visible to both the place you use the calendar heatmap link and the server that generates the calendar heatmap. If your pomo history has some important information you don't want to leak, you should not use this service. The code is open source as you see, I run it on my server. If you don't trust me, you can self-host pomotodo-stats. Download the source code and build it yourself. Or use [the docker](https://hub.docker.com/repository/docker/wintercicada/pomotodo-stats) to deploy it.

## Todo
- [ ] Custom heatmap color and date range
- [ ] Add most used tags view



## Credits
[blurfx/calendar-heatmap](https://github.com/blurfx/calendar-heatmap)

[arran4/golang-ical](https://github.com/arran4/golang-ical)


