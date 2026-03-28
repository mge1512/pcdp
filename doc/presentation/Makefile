default:	pcd-introduction.pdf 

pictures=	pcd-icon-green.png pcd-logo-green.png 								\
		pcd-workflow-pcd.png pcd-workflow.png pcd-workflow-today.png pcd-resolution.png 			\
		pcd-spec-sections.png pcd-audit-bundle.png pcd-translation-phases.png pcd-mcp-server.png

pcd-introduction.md: pcd-intro-00-header.md pcd-intro-10-overview.md pcd-intro-50-technical.md pcd-intro-90-appendix.md
	cat $^ > $@

pcd-introduction.pdf: pcd-introduction.tex $(pictures)

%.pdf: %.tex
	pdflatex $< ; pdflatex $<

%.tex: %.md
	pandoc --to=beamer -V classoption:10pt --slide-level=2 --standalone --output=$@ $+

%.docx: %.md
	pandoc --to=docx --slide-level=2 --standalone --output=$@ $+

pcd-icon-green.png: pcd-icon-green.svg 
	magick -background none -size 50x50 $< $@

%.png: %.svg
	magick -background none  $< $@

%.svg: %.pikchr
	pikchr --svg-only $< | sed 's/<text /<text font-family="DejaVu Sans, Liberation Sans, sans-serif" /g' > $@

distclean:
	rm pcd-introduction.md pcd-introduction.aux pcd-introduction.log pcd-introduction.nav pcd-introduction.pdf pcd-introduction.snm pcd-introduction.toc pcd-introduction.vrb *.svg *.png

clean:
	rm pcd-introduction.md pcd-introduction.aux pcd-introduction.log pcd-introduction.nav pcd-introduction.tex pcd-introduction.snm pcd-introduction.toc pcd-introduction.vrb *.svg *.png

view:	pcd-introduction.pdf
	evince -s pcd-introduction.pdf

