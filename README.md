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

#### Description

Users usually make typos or simply do not know how to spell station names correctly. That's why we need to implement an
approximate search for the station names.

Let's define the closest to a sample word as:

- It has the longest common substring. Length of the substring has to be more than 0
- It has the same set of letters or, at least, this set has minimal difference with sample's set

This definition allows us to filter completely unsuitable words - if a word has no matching substrings with the sample,
it means they have 0 common letters and then can not be called 'close'.

At the same time our definition explicitly marks the closest word not only a word with the closest letters set match
but also uses information about letters order.

The search steps are:

1. First, for each character in the sample word, we count how many times it appears.
   We also make a list of all 'end parts' (suffixes) of the word that start with the same character.
2. Then, for each word we're comparing to the sample, we do the same - we list out 'end parts' (suffixes) and
   check them against the 'end parts' list we made previously for the sample. We're specifically trying to match the
   beginnings of the 'end parts'. We keep track of the longest match we find.
3. We also tally up the characters used in each word and compare these 'character sets' between the sample and each
   word.
4. At the end, we choose the word that had the smallest difference in the 'character sets' and the longest match in
   'end parts'. This is our best match.

#### Complexity

The time complexity of this algorithm can be determined by examining the operations performed:

1. Building RuneStat for the sample string requires scanning each character for a total of O(n) operations,
   where n is the length of the sample string.
2. The next loop involves iterating over each word in the search list. The worst-case scenario for each word is when
   it equals the sample—the loops would result in O(m^2) operations, where m is the size of the word.
   As this is performed for each word in the searchList, the worst-case time complexity becomes O(k * m^2),
   where k is the size of the searchList.
3. To find the minimum word, we need to scan all the words stored, this procedure will take O(k).

Given this, the overall time complexity is dominated by the second step, resulting in O(m^2 * k) in the worst-case
scenario.

Space complexity depends on the storage of the intermediate RuneStats and Words:

1. Storage of the RuneStat objects for the sample string is O(n).
2. Within the loop, the creation of the words slice can lead to a space complexity of O(k).

Therefore, the worst-case space complexity of the algorithm is O(n + k), considering the storage of RuneStat objects
for n sample characters and the storage for k words in the search list

// TODO: compare with others:
https://github.com/lithammer/fuzzysearch ,
https://github.com/schollz/closestmatch ,
https://github.com/antzucaro/matchr

