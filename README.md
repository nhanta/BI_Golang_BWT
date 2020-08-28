# Applied Concurrency Technique in MultipleApproximate Pattern Matching Problem with Burrows-Wheeler Transform
#
**1 Introduction**

We focused on designing algorithms on the basis of the Burrows-Wheeler Transform algorithm improvement to achieve the highest efficiency when deploying and taking advantage of Golang's strengths. We hope that this is an approach for testing the new programming language that has advantages in speed, popularity, and ease of installation on universal computing systems for the bioinformatics problem in general, and the problem of sequence alignment in particular.

**2 Data**

_We retrieved the [raw sequences of SARS-CoV-2](https://sra-pub-sars-cov2.s3.amazonaws.com/sra-src/SRR12338312/KPCOVID-345_S81_L001_R1_001.fastq.gz.1) published on July 28, 2020 by KwaZulu-Natal Research Innovation and Sequencing Platform from the Sequence Read Archive (SRA). The FASTQ file includes 436.610 paired-end reads. The FASTQ file  was converted to the fasta file (named **Sra\_SARs\_CoV\_2.fasta**) by the tool FASTQ to FASTA converter on Galaxy Version 1.1.5.

_[The genome assembly of SARS-CoV-2](https://www.ncbi.nlm.nih.gov/nuccore/NC_045512.2) published by Fan Wu et al. (2020), which is 24748 bp long was used as the reference genome for alignment. The reference genome file was renamed to **Ref\_SARs\_CoV\_2.fa**.

**2 Preparation**

_The code is run with Python 3.7.1. Pandas package need to be installed before running the code._

Using the data from _Prepared\_data.rar_.

Filtering the data with p value, then put it into the models for training. We have also saved the results in _Export\_data.rar._

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

