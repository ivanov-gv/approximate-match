# approximate-match

Let's assume we have a list of the words, and we want to find the closest match to another one - sample.
For example:  

aaaaa - our sample 

List:
aaaab  
aaaac  
abcde  
bbbbb  


Let's define the closest to sample word as:
1. It has the longest common substring, but with length of this substring has to be more than 0
2. It has the same set of letters or, at least, this set of letters has minimal difference with sample's set

As in our example before:  
aaaab  - has substring "aaaa", 1 less of letters a, 1 more of letters b. So 4 - is a lenght of a common substring, 2 - set difference  
aaaac  - the same as before - 4 and 2
abcde  - length - 1 , 8 - difference (-4 'a', + 'b', + 'c', + 'd', + 'e' )  
bbbbb  - 0 common letters, difference is 10

According to our definition "aaaab", "aaaac" and "abcde" are our candidates, but "bbbbb" is not.

Let's use it for searching another closest match.
Our next sample - Ella
List: Adele Elaine Elizabeth Harriet Ingrid Michelle Ella

Output:
Ella      - exact match, no differences
Michelle  - has the same substring 'ell', but has 5+ more letters and doesn't have 1 letter 'a' 
Adele  
Elaine  
Elizabeth  
Harriet  

and no word "Ingrid" in output list, because it doesn't match on any letter.



