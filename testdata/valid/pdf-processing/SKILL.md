---
name: pdf-processing
description: Extract text and tables from PDF files, fill forms, merge documents. Use when working with PDF documents or when the user mentions PDFs, forms, or document extraction.
license: Apache-2.0
compatibility: Designed for Claude Code (or similar products)
metadata:
  author: example-org
  version: "1.0"
---

# PDF Processing

## When to use this skill

Use this skill when the user needs to work with PDF files, including:
- Extracting text from PDFs
- Extracting tables from PDFs
- Filling PDF forms
- Merging multiple PDFs

## How to extract text

1. Use pdfplumber for text extraction
2. Handle multi-page documents
3. Return structured text

## How to fill forms

1. Identify form fields
2. Map user data to fields
3. Generate filled PDF
