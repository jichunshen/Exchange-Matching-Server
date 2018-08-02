library(dplyr)
library(tidyverse)
iris=as.tibble(iris)
iris
tbl_df(cars)
cars
a=4 %>% .$+$2
a
iris %>%
  group_by(Species)%>%
  summarise(avg=mean(Sepal.Width))%>%
  arrange(avg)
  

tapply(seq(1,4),$*$4)
x <- list(a = 1, b = 1:3, c = 10:100) 


x <- list(a = 1:10, beta = exp(-3:3), logic = c(TRUE,FALSE,FALSE,TRUE))
# compute the list mean for each list element
lapply(x,mean)
# median and quartiles for each list element
lapply(x, quantile, probs = 1:5/6)
sapply(x, quantile)
i39 <- sapply(3:9, seq) %>%
sapply(scale)%>%
  print
library(MASS)
apply(iris[,1:4],2,function(x){mean(x)^2+3})

3.172/sqrt(2.24*6.97)

lapply(iris[,1:4],function(x){(x-mean(x))}) %>%
  as.data.frame%>%
  apply(.,2,mean)


827/sqrt(98775*19.1)



as.tibble(iris) %>%
   filter(Species=="setosa")

 1-34/182
 
 
r$a
apply(iris[,1:4],2,sd)

a=seq(5) %>%
  (function(x){x^2+4-6}) %>%
  abs %>%
  sqrt
a
z=iris %>%
  filter(., Sepal.Width %in% seq(3,4,0.1))
summary(iris)


a=3 %>%
  sqrt 
library(rbokeh)


figure(iris) %>%
  ly_points(Sepal.Length, Sepal.Width, color=Species)




 %>%
   %>%
  %>%
  dyOptions( dySeries(dySeries(dygraph(lungDeaths),"mdeaths", label = "Male"),"fdeaths", label = "Female"),stackedGraph = TRUE) %>%
  dyRangeSelector(height = 20)


<- function(variables) {
  
}(iris)
z=filter(iris, Species=="virginica") 
z


x=tibble(id=c(1,1,1,2,2,2,3,3,3),e=as.factor(c(1,2,3,1,2,3,1,2,3)),b=c(99,89,72,35,22,66,99,99,98))
x
x$id=x$id %>%
  as.factor(.)
spread(x,id,e)
?spread



messy <- data.frame(
  name = c("LI", "PW", "Gregory"),
  exam1 = c(67, 80, 64),
  exam2 = c(56, 90, 50),
  exam3 = c(30,20,40),
  exam4 = c(30,20,40),
  exam5 = c(30,20,40),
  exam6 = c(30,20,40)
  
)
messy
messy=messy %>%
  gather(exam, score, exam1:exam6) %>%
  arrange(name) %>%
  rename(allexams=exam)

spread(messy,allexams,3)
as.tibble(distinct(iris))
sample_n(iris, 10, replace = FALSE)
install.packages("data.table")
library(data.table)
filter(iris, Species %like% "er" )
