# Makefile for Quine A-Life Track Paper

# Output directory
OUT_DIR := output
# Source files are in the sections directory
SRC_DIR := sections
META := metadata.yaml

# Source files (order matters)
SECTIONS := $(SRC_DIR)/00-abstract.md \
            $(SRC_DIR)/01-introduction.md \
            $(SRC_DIR)/02-background.md \
            $(SRC_DIR)/03-anatomy.md \
            $(SRC_DIR)/04-experiments.md \
            $(SRC_DIR)/05-limitations.md \
            $(SRC_DIR)/06-future.md \
            $(SRC_DIR)/07-conclusion.md \
            $(SRC_DIR)/08-references.md \
            $(SRC_DIR)/appendix-a.md \
            $(SRC_DIR)/appendix-b.md \
            $(SRC_DIR)/appendix-c.md

# Figure files
FIGURES := $(wildcard figures/*.svg)

# Combined markdown file
FULL_MD := $(OUT_DIR)/full_paper.md

# Output files
PDF_NAME := Quine_POSIX_as_Physics_for_Emergent_Artificial_Life.pdf
PDF := $(OUT_DIR)/$(PDF_NAME)
TEX := $(OUT_DIR)/paper.tex
HTML := $(OUT_DIR)/paper.html
DOCX := $(OUT_DIR)/paper.docx

# ArXiv submission directory
ARXIV_DIR := arxiv_submission
ARXIV_TEX := $(ARXIV_DIR)/main.tex
ARXIV_ZIP := $(ARXIV_DIR).tar.gz

# Pandoc options (use tectonic for PDF generation)
PANDOC_OPTS := --standalone \
               --from markdown \
               --pdf-engine=tectonic \
               --resource-path=.:sections \
               --variable fontsize=11pt \
               --highlight-style tango

.PHONY: all clean pdf tex html docx setup arxiv arxiv-clean

all: pdf arxiv

setup:
	mkdir -p $(OUT_DIR)

# Combine metadata and all sections into one markdown file
# Use awk to ensure blank line between each file (for proper markdown section separation)
$(FULL_MD): $(META) $(SECTIONS) | setup
	awk 'FNR==1 && NR!=1 {print ""} {print}' $(META) $(SECTIONS) > $(FULL_MD)

# Generate PDF
pdf: setup $(FULL_MD)
	pandoc $(FULL_MD) -o $(PDF) $(PANDOC_OPTS)
	@echo "PDF generated at $(PDF)"

# Generate TeX
tex: $(FULL_MD)
	pandoc $(FULL_MD) -o $(TEX) --standalone --from markdown --to latex --resource-path=.:sections
	@echo "TeX generated at $(TEX)"

# Generate HTML
html: $(FULL_MD)
	pandoc $(FULL_MD) -o $(HTML) --standalone --resource-path=.:sections
	@echo "HTML generated at $(HTML)"

# Generate DOCX
docx: $(FULL_MD)
	pandoc $(FULL_MD) -o $(DOCX) --resource-path=.:sections
	@echo "DOCX generated at $(DOCX)"

# Generate ArXiv submission package
.PHONY: arxiv
arxiv: $(FULL_MD)
	@echo "=== Building ArXiv submission package ==="
	rm -rf $(ARXIV_DIR) $(ARXIV_ZIP)
	mkdir -p $(ARXIV_DIR)/figures
	
	@# Convert SVG figures to PDF (arXiv requires PDF/PNG/JPG, not SVG)
	@echo "Converting figures from SVG to PDF..."
	@for svg in figures/*.svg; do \
		pdf="$(ARXIV_DIR)/figures/$$(basename "$$svg" .svg).pdf"; \
		echo "  $$svg -> $$pdf"; \
		rsvg-convert -f pdf -o "$$pdf" "$$svg"; \
	done
	
	@# Generate TeX with ArXiv-compatible settings (pdflatex, no custom fonts)
	@echo "Generating LaTeX..."
	pandoc $(FULL_MD) -o $(ARXIV_TEX) \
		--standalone \
		--from markdown \
		--to latex \
		--pdf-engine=pdflatex \
		--variable documentclass=article \
		--variable fontsize=11pt \
		--variable geometry=margin=1in \
		--variable mainfont= \
		--variable monofont=
	
	@# Post-process for ArXiv compatibility
	@echo "Post-processing TeX for ArXiv compatibility..."
	@# Replace .svg with .pdf in figure references
	sed -i.bak 's/\.svg/.pdf/g' $(ARXIV_TEX)
	@# Replace \includesvg with \includegraphics (svg package won't work with PDF files)
	sed -i.bak 's/\\includesvg/\\includegraphics/g' $(ARXIV_TEX)
	@# Remove svg package if present
	sed -i.bak '/\\usepackage{svg}/d' $(ARXIV_TEX)
	@# Remove lmodern (use default Computer Modern which arXiv definitely has)
	sed -i.bak '/\\usepackage{lmodern}/d' $(ARXIV_TEX)
	@# Remove entire microtype block (IfFileExists check and its contents)
	sed -i.bak '/IfFileExists{microtype.sty}/,/^}{}/d' $(ARXIV_TEX)
	rm -f $(ARXIV_TEX).bak
	
	@# Remove test-generated PDF if exists
	rm -f $(ARXIV_DIR)/main.pdf $(ARXIV_DIR)/*.aux $(ARXIV_DIR)/*.log $(ARXIV_DIR)/*.out
	
	@# Create tarball
	@echo "Creating submission tarball..."
	tar -czvf $(ARXIV_ZIP) -C $(ARXIV_DIR) .
	
	@echo ""
	@echo "=== ArXiv submission package ready ==="
	@echo "Directory: $(ARXIV_DIR)/"
	@echo "Tarball:   $(ARXIV_ZIP)"
	@echo ""
	@echo "Package size:"
	@du -sh $(ARXIV_DIR)/ $(ARXIV_ZIP)
	@echo ""
	@echo "Contents:"
	@ls -lh $(ARXIV_DIR)/
	@echo ""
	@echo "Next steps:"
	@echo "  1. Test compile: cd $(ARXIV_DIR) && pdflatex main.tex"
	@echo "  2. Upload $(ARXIV_ZIP) to arxiv.org"

arxiv-clean:
	rm -rf $(ARXIV_DIR) $(ARXIV_ZIP)

clean:
	rm -f $(FULL_MD) $(PDF) $(TEX) $(HTML) $(DOCX)
	rm -rf $(ARXIV_DIR) $(ARXIV_ZIP)
	rm -rf $(OUT_DIR)
