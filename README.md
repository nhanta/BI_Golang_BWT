# Applied Concurrency Technique in Multiple Approximate Pattern Matching Problem with Burrows-Wheeler Transform
#
**1 Introduction**

We focused on designing algorithms on the basis of the Burrows-Wheeler Transform algorithm improvement to achieve the highest efficiency when deploying and taking advantage of Golang's strengths. We hope that this is an approach for testing the new programming language that has advantages in speed, popularity, and ease of installation on universal computing systems for the bioinformatics problem in general, and the problem of sequence alignment in particular.

**2 Data**

We retrieved the [raw sequences of SARS-CoV-2](https://sra-pub-sars-cov2.s3.amazonaws.com/sra-src/SRR12338312/KPCOVID-345_S81_L001_R1_001.fastq.gz.1) published on July 28, 2020 by KwaZulu-Natal Research Innovation and Sequencing Platform from the Sequence Read Archive (SRA). The FASTQ file includes 436.610 paired-end reads. The FASTQ file  was converted to the fasta file (named **Sra_SARs_CoV_2.fasta**) by the tool FASTQ to FASTA converter on Galaxy Version 1.1.5.

[The genome assembly of SARS-CoV-2](https://www.ncbi.nlm.nih.gov/nuccore/NC_045512.2) published by Fan Wu et al. (2020), which is 24748 bp long was used as the reference genome for alignment. The reference genome file was renamed to **Ref_SARs_CoV_2.fa**.

**3 Preparation**

Download the project and organize the folders as follows:
```
MAPMBWT
|-- Packages            
|   |-- CheckPoint                       Create_CheckPoint Arrays_ with different parameters k.
|   `-- CheckPointArrays.go      
|   |-- ConverttoByte                    Convert an integer to a byte.
|   `-- InttoByte.go
|   |-- MemUsage                         Print memory usage.
|   `-- PrintMemUsage.go
|   |-- PartialSuffixArrays              Create _Partial Suffix Arrays_ with different parameters c.
|   `-- partialsuffixarrays.go
|   |-- ReadFiles                        Write and read files.
|   `-- readlines.go
|   `-- readwords.go
|   `-- rwjson.go
|   `-- rwstringtotext.go
|   |-- ReverseSeq                       Convert a sequence to its complementary. 
|   `-- ReverseSeq.go
|   |--TexttoBWT                         Creating _Burrows-Wheeler Transform._
|   `-- texttobwt.go
|-- index.go                             Export _CheckPoint Arrays_, _Partial Suffix Arrays_, and _Burrows-Wheeler Transform._
|-- main.go                              Multiple Approximate Pattern Matching Algorithm.
|-- compare.go                           Classify match sequences according to difference threshold.
|-- align.go                             Find position of SNPs.
|-- Sra_SARs_CoV_2.fasta                 SRA file.
|-- Ref_SARs_CoV_2.fa                    Reference genome file.

```

**3 Code-Running Steps**

Step 1: _go run index.go_
Step 2: _go run main.go_
Step 3: _go run compare.go_
Step 4: _go run align.go_
