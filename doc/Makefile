default:	pcdp-introduction.pdf 

pictures:	pcdp-workflow.png pcdp-resolution.png pcdp-icon-green.png pcdp-logo-green.png pcdp-icon-green.png pcdp-workflow-pcdp.png pcdp-workflow.png pcdp-workflow-today.png

pcdp-introduction.md: pcdp-intro-00-header.md pcdp-intro-10-overview.md pcdp-intro-50-technical.md pcdp-intro-90-appendix.md
	cat $^ > $@

%.pdf: %.tex
	pdflatex $< ; pdflatex $<

%.tex: %.md
	pandoc --to=beamer -V classoption:10pt --slide-level=2 --standalone --output=$@ $+

%.docx: %.md
	pandoc --to=docx --slide-level=2 --standalone --output=$@ $+

pcdp-icon-green.png: pcdp-icon-green.svg 
	magick -background none -size 50x50 $< $@

%.png: %.svg
	magick -background none  $< $@

%.svg: %.pikchr
	pikchr --svg-only $< | sed 's/<text /<text font-family="DejaVu Sans, Liberation Sans, sans-serif" /g' > $@

clean:
	rm pcdp-introduction.md pcdp-introduction.aux pcdp-introduction.log pcdp-introduction.nav pcdp-introduction.pdf pcdp-introduction.snm pcdp-introduction.toc pcdp-introduction.vrb *.svg *.png

view:	pcdp-introduction.pdf
	evince -s pcdp-introduction.pdf

