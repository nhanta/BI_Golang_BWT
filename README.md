# Applied Concurrency Technique in MultipleApproximate Pattern Matching Problem with Burrows-Wheeler Transform
#
**1 Introduction**

We focused on designing algorithms on the basis of the Burrows-Wheeler Transform algorithm improvement to achieve the highest efficiency when deploying and taking advantage of Golang's strengths. We hope that this is an approach for testing the new programming language that has advantages in speed, popularity, and ease of installation on universal computing systems for the bioinformatics problem in general, and the problem of sequence alignment in particular.

**2 Data**

We retrieved the [raw sequences of SARS-CoV-2](https://sra-pub-sars-cov2.s3.amazonaws.com/sra-src/SRR12338312/KPCOVID-345_S81_L001_R1_001.fastq.gz.1) published on July 28, 2020 by KwaZulu-Natal Research Innovation and Sequencing Platform from the Sequence Read Archive (SRA). The FASTQ file includes 436.610 paired-end reads. The FASTQ file  was converted to the fasta file (named **Sra\_SARs\_CoV\_2.fasta**) by the tool FASTQ to FASTA converter on Galaxy Version 1.1.5.

[The genome assembly of SARS-CoV-2](https://www.ncbi.nlm.nih.gov/nuccore/NC_045512.2) published by Fan Wu et al. (2020), which is 24748 bp long was used as the reference genome for alignment. The reference genome file was renamed to **Ref\_SARs\_CoV\_2.fa**.

**2 Preparation**

Download the project and organize the folders as follows:
```
MAPMBWT
|-- README.md                  This README file.
|-- run-bwamem                 *Entry script* for the entire mapping pipeline.
|-- bwa                        *BWA binary*
|-- k8                         Interpretor for *.js scripts.
|-- bwa-postalt.js             Post-process alignments to ALT contigs/decoys/HLA genes.
|-- htsbox                     Used by run-bwamem for shuffling BAMs and BAM=>FASTQ.
|-- samblaster                 MarkDuplicates for reads from the same library. v0.1.20
|-- samtools                   SAMtools for sorting and SAM=>BAM conversion. v1.1
|-- seqtk                      For FASTQ manipulation.
|-- trimadap                   Trim Illumina PE sequencing adapters.
|
|-- run-gen-ref                *Entry script* for generating human reference genomes.
|-- resource-GRCh38            Resources for generating GRCh38
|   |-- hs38DH-extra.fa        Decoy and HLA gene sequences. Used by run-gen-ref.
|   `-- hs38DH.fa.alt          ALT-to-GRCh38 alignment. Used by run-gen-ref.
|
|-- run-HLA                    HLA typing for sequences extracted by bwa-postalt.js.
|-- typeHLA.sh                 Type one HLA-gene. Called by run-HLA.
|-- typeHLA.js                 HLA typing from exon-to-contig alignment. Used by typeHLA.sh.
|-- typeHLA-selctg.js          Select contigs overlapping HLA exons. Used by typeHLA.sh.
|-- fermi2.pl                  Fermi2 wrapper. Used by typeHLA.sh for de novo assembly.
|-- fermi2                     Fermi2 binary. Used by fermi2.pl.
|-- ropebwt2                   RopeBWT2 binary. Used by fermi2.pl.
|-- resource-human-HLA         Resources for HLA typing
|   |-- HLA-ALT-exons.bed      Exonic regions of HLA ALT contigs. Used by typeHLA.sh.
|   |-- HLA-CDS.fa             CDS of HLA-{A,B,C,DQA1,DQB1,DRB1} genes from IMGT/HLA-3.18.0.
|   |-- HLA-ALT-type.txt       HLA types for each HLA ALT contig. Not used.
|   `-- HLA-ALT-idx            BWA indices of each HLA ALT contig. Used by typeHLA.sh
|       `-- (...)
|
`-- doc                        BWA documentations
    |-- bwa.1                  Manpage
    |-- NEWS.md                Release Notes
    |-- README.md              GitHub README page
    `-- README-alt.md          Documentation for ALT mapping
```


**3 Training Model**

Python code: **BI\_Rice\_SNPs.py**

_The code is run with Python 3.7.1. Pandas package and scikit-learn package need to be installed before running the code._

The data from _Export\_data.rar_ is used to train the model.


3.1 Grid Search with cross-validation: [sklearn.model\_selection.GridSearchCV](https://scikit-learn.org/stable/modules/generated/sklearn.model_selection.GridSearchCV.html)

Link: https://scikit-learn.org/stable/modules/generated/sklearn.model\_selection.GridSearchCV.html

Used with random forest regression and support vector regression.

3.2 Radom forest regression: [sklearn.ensemble.RandomForestRegressor](https://scikit-learn.org/stable/modules/generated/sklearn.ensemble.RandomForestRegressor.html)

Link: https://scikit-learn.org/stable/modules/generated/sklearn.ensemble.RandomForestRegressor.html

3.3 Support vector regression: [sklearn.svm.SVR](https://scikit-learn.org/stable/modules/generated/sklearn.svm.SVR.html)

Link: https://scikit-learn.org/stable/modules/generated/sklearn.svm.SVR.html

3.4 Lasso with cross-validation: [sklearn.linear\_model.LassoCV](https://scikit-learn.org/stable/modules/generated/sklearn.linear_model.LassoCV.html)

Link: https://scikit-learn.org/stable/modules/generated/sklearn.linear\_model.LassoCV.html

3.5 Multi-task Lasso with cross-validation: [sklearn.linear\_model.MultiTaskLassoCV](https://scikit-learn.org/stable/modules/generated/sklearn.linear_model.MultiTaskLassoCV.html)

Link: https://scikit-learn.org/stable/modules/generated/sklearn.linear\_model.MultiTaskLassoCV.html

3.6 Elastic Net with cross-validation: [sklearn.linear\_model.ElasticNetCV](https://scikit-learn.org/stable/modules/generated/sklearn.linear_model.ElasticNetCV.html)

Link: https://scikit-learn.org/stable/modules/generated/sklearn.linear\_model.ElasticNetCV.html

3.7 Multi-task Elastic Net with cross-validation: [sklearn.linear\_model.MultiTaskElasticNetCV](https://scikit-learn.org/stable/modules/generated/sklearn.linear_model.MultiTaskElasticNetCV.html)

Link: https://scikit-learn.org/stable/modules/generated/sklearn.linear\_model.MultiTaskElasticNetCV.html

The results have been compressed into the file _Results.rar._

