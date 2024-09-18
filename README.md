# gorbit

My 6-year-old son wanted one of those toy Sun-Earth-Moon orbital models, but shopping online presented me with products that were inaccuate, low quality, unattractive or exorbitantly expensive. I had been meaning to learn Golang for the two years my current job has exposed me to it. So I decided this was an opportunity to learn with an objective in mind, and I named it gorbit (for Go-Orbit)

As my first foray into involved Go programming, one might notice my first attempt can be rather convoluted. In my defense, I was also learning astrophysics through Google. At least I got around to breaking up an even-more-unwieldy single main.go file into objects that were more meaningful. It uses the G3N ("gen") OpenGL 3D Game Engine for Go, which enabled my geometric design without needing to learn hardware system calls


NOTES

Intuition based on Northern Hemisphere, for the current year of 2024 at the time of this writing.

Vernal Equinox 2024-03-19

Summer Solstice 2024-06-20 00:00 (Initial Earth position) Epoch Offset 1718841600 seconds. Also happens to be when moon is approximately at Greenwich, England meridian, luckily making initial positioning slightly easier

Autumnal Equinox 2024-09-22

Winter Solstice 2024-12-21


ASSUMPTIONS

No Euclidean elliptical orbits for now

Lunar plane rotational-Y orientation stays absolutely fixed beyond sun, i.e., it is static to the Earth's tilt axis which results in our seasons, and not to the revolution around the Sun


TO DO

Forward backward direction of animation speed

More guidelines that show angles, periods of events

Real time

Input fields for setting time manually

Reset button back to initial solstice, or to actual time

Scale distances and sizes correctly

Animating camera reorientation and re-targeting

Render shadow for moon to cast for solar eclipses

Lunar plane angle oscillates over ~18 years

Change simplified circular orbits to actually elliptical


REFERENCES

https://www.timeanddate.com/moon

https://code.visualstudio.com/docs/cpp/config-mingw

https://github.com/g3n/engine

https://pkg.go.dev/github.com/g3n/engine
