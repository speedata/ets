# ETS Experimental typesetting system

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
-- The document d is the central point for all
function CreateFontVlist(d)
    local face, msg = d.loadFace("fonts/CrimsonPro-Regular.ttf")
    if not face then
        print(msg)
        os.exit(-1)
    end
    local fnt = d.createFont(face,document.sp("12pt"))

    local lang_en, msg = d.loadpattern("hyphenationpatterns/hyph-en-us.pat.txt")
    if not lang_en then
        print(msg)
        os.exit(-1)
    end

    lang_en.name = "en"

    local lang = node.new("lang")
    lang.lang = lang_en

    local tbl = fnt.shape("the quick brown fox jumps over the lazy dog")


    local head, cur = lang, lang
    for _, glyph in ipairs(tbl) do
        if glyph.glyph == 32 then
            local glu = node.new("glue")
            glu.width = fnt.space
            head = node.insertafter(head,cur,glu)
            cur = glu
        else
            local g = node.new("glyph")
            g.width = glyph.advance
            g.codepoint = glyph.codepoint
            g.font = glyph.font
            head = node.insertafter(head,cur,g)
            cur = g
        end
    end

    local hlist = node.hpack(head)

    local param = {
        hsize = document.sp("200pt"),
        lineheight = document.sp("12pt"),
    }

    local vl = node.simplelinebreak(hlist,param)
    return vl
end

local function CreateImageVlist(d)
    local imgfile = d.loadimagefile("img/ocean.pdf")
    local image = d.createimage(imgfile)


    local imagenode = node.new("image")
    imagenode.img = image
    imagenode.width = document.sp("4cm")
    imagenode.height = document.sp("3cm")
    local vlist = node.new("vlist")
    vlist.list = imagenode
    return vlist
end


local d = document.new("out.pdf")

local fontVL = CreateFontVlist(d)
local imageVL = CreateImageVlist(d)

d.outputat(document.sp("4cm"),document.sp("27cm"),fontVL)
d.outputat(document.sp("4cm"),document.sp("26cm"),imageVL)
d.currentpage().shipout()


local ok, msg = d.finish()
if not ok then
    print(msg)
    os.exit(-1)
end
```

