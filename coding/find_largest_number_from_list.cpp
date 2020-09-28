#include <iostream>
#include <vector>
#include <string> 
#include <algorithm> 

using namespace std;

int compareFn(int a, int b) {
    string x = to_string(a);  
    string y = to_string(b); 
    
    string xy = x.append(y);
    string yx = y.append(x);
    
    return xy.compare(yx) > 0 ? 1: 0;
}

int main()
{
   vector<int> a;
   a.push_back(1);
   a.push_back(19);
   a.push_back(2);
   
   // sort the array
   sort(a.begin(), a.end(), compareFn);
   
   // print the sorted array
   for (int i=0; i<a.size(); i++) {
       cout << a[i];
   }
   return 0;
}