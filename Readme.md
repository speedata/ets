# ets Experimental typesetting system

This software repository contains a Lua frontend for the typesetting library [“Boxes and Glue”](https://github.com/speedata/boxesandglue) which is an algorithmic typesetting machine in the spirits of TeX.


## Running the software

Download the software ([get the latest release here](https://github.com/speedata/ets/releases)), unzip the archive and call

    bin/ets myscript.lua

on the command line.

## Build

You just need Go installed on your system, clone the repository and run

    go build -o bin/ets github.com/speedata/ets/ets/ets



## Status

This software is more or less a demo of the architecture and not usable for any serious purpose yet.

Feedback welcome!


## Contact and information

License: AGPL v3<br>
Contact: Patrick Gundlach <gundlach@speedata.de>



## Sample code

```lua
-- This is still a very vey early preview and will probably not work when
-- you experiment with the code.
--
-- Two functions for fonts and for images just to show how to
-- create objects and output them.

document.info("Reading myfile.lua")

local ok, face, msg, fnt, lang_en, imgfile, image

function CreateFontVlist(d)
    local lang = node.new("lang")
    lang.lang = lang_en

    local tbl = fnt.shape([[In olden times when wishing still helped one, there lived a king whose daughters
were all beautiful; and the youngest was so beautiful that the sun itself, which
has seen so much, was astonished whenever it shone in her face.
Close by the king's castle lay a great dark forest, and under an old lime-tree in the forest
was a well, and when the day was very warm, the king's child went out into the
forest and sat down by the side of the cool fountain; and when she was bored she
took a golden ball, and threw it up on high and caught it; and this ball was her
favorite plaything.]])


    local head, cur = lang, lang
    for _, glyph in ipairs(tbl) do
        if glyph.glyph == 32 then
            local glu = node.new("glue")
            glu.width = fnt.space
            glu.stretch = fnt.stretch
            glu.shrink = fnt.shrink
            head = node.insertafter(head,cur,glu)
            cur = glu
        else
            local g = node.new("glyph")
            g.width = glyph.advance
            g.codepoint = glyph.codepoint
            g.components = glyph.components
            g.font = glyph.font
            g.hyphenate = glyph.hyphenate
            head = node.insertafter(head,cur,g)
            cur = g
        end
    end
    d.hyphenate(head)
    node.append_lineend(cur)

    local param = {
        hsize = document.sp("134pt"),
        lineheight = document.sp("12pt"),
    }

    local vl = node.linebreak(head,param)
    return vl
end

local function CreateImageVlist(d)
    image = d.createimage(imgfile)
    local imagenode = node.new("image")
    imagenode.img = image
    imagenode.width = document.sp("4cm")
    imagenode.height = document.sp("3cm")
    local vlist = node.new("vlist")
    vlist.list = imagenode
    return vlist
end

-- The document d is the most important item here.
local d = document.new("out.pdf")

face, msg = d.loadFace("fonts/CrimsonPro-Regular.ttf")
if not face then
    print(msg)
    os.exit(-1)
end

fnt = d.createFont(face,document.sp("12pt"))

lang_en, msg = d.loadpattern("hyphenationpatterns/hyph-en-us.pat.txt")
if not lang_en then
    print(msg)
    os.exit(-1)
end

lang_en.name = "en"
imgfile = d.loadimagefile("img/ocean.pdf")

local fontVL = CreateFontVlist(d)
local imageVL = CreateImageVlist(d)
d.outputat(document.sp("4cm"),document.sp("27cm"),fontVL)
d.outputat(document.sp("12cm"),document.sp("27cm"),imageVL)
d.currentpage().shipout()

ok, msg = d.finish()
if not ok then
    print(msg)
    os.exit(-1)
end

document.info("Reading myfile.lua...done")
```

