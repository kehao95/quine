PAPER_TEMPLATE_ROOT := ../../templates/arxiv-acmart

TEMPLATE := $(PAPER_TEMPLATE_ROOT)/templates/acmart.latex
SUPPORT_DIR := $(PAPER_TEMPLATE_ROOT)/support/acmart

SECTIONS := sections/01-introduction.md \
            sections/02-problem-related-work.md \
            sections/03-operational-assay.md \
            sections/04-experimental-design.md \
            sections/05-core-results.md \
            sections/06-controls-robustness.md \
            sections/07-discussion.md \
            sections/08-limitations-conclusion.md \
            sections/09-artifact-appendix.md

PDF_NAME := The_Emergence_of_Self_Replication.pdf
TEX := output/paper.tex

include $(PAPER_TEMPLATE_ROOT)/Makefile

# The installed Pandoc 3.7 binary is paired with a system data package that
# omits this optional partial. This paper does not use document-level PDF
# metadata, so provide a local no-op compatibility overlay.
PANDOC_PDF_OPTS += --data-dir=$(CURDIR)/pandoc-data
PANDOC_LATEX_OPTS += --data-dir=$(CURDIR)/pandoc-data
