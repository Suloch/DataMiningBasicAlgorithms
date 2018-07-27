
package main

import (
  "fmt"
  "bufio"
  "os"
  "strings"
  "strconv"
)

//structure for holding rows of a table that is an item set and support count
type row struct{
  item_set []int
  sup_count int
}

type table struct{
  rows []row
}


func main(){
  var L []table
  transactions := readFile("data")
  const minimum_support = 2

  L = append(L, createL1(transactions))

  printTable(L[0])

  for i:=1; len(L[i-1].rows) > 0; i++{
    L = append(L, generateC(L[i-1]))

    for j,row1 := range L[i].rows{
      L[i].rows[j].sup_count = countSupport(row1.item_set, transactions)
    }

    //delete rows with support less than minimum support
    for j:=0; j< len(L[i].rows);{
      if L[i].rows[j].sup_count < minimum_support{
        L[i].rows[j] = L[i].rows[len(L[i].rows)-1]
        L[i].rows = L[i].rows[:len(L[i].rows)-1]
      }  else{
        j++;
      }
    }
    printTable(L[i])
  }


//  A, B, AUB := getInput()

}

//function to get input from the user
func getInput() ([]int, []int, []int){

  var A, B, AUB []int


  return A,B,AUB
}

//function to count the support of an item set
func countSupport(item_set []int, transactions [][]int) int{

  present := false
  sup_count := 0
  for _,transaction := range transactions{
    for _,item1 := range item_set{
      present = false
      for _,item2 := range transaction{
        if item1 == item2{
          present = true
          break
        }
      }
      if !present{
        break
      }
    }
    if present{
      sup_count++
    }
  }
  return sup_count
}

//function to generate the next C table from given L table
func generateC(L table) table{
  var C table

  for _,row1 := range L.rows{
    for _,row2 := range L.rows{
      row3, merged := mergeRow(row1, row2)
      if merged{
        present := false
        for _, row4 := range C.rows{
          var checked = make([]bool, len(row4.item_set))
          for i,item := range row4.item_set{
            checked[i] = false
            for _,item1 := range row3.item_set{
              if item1 == item{
                checked[i] = true
                break;
              }
            }
          }
          for i:=0; i < len(checked); i++{
            present = false;
            if !checked[i]{
              break;
            }
            present = true
          }

          if present{
            break;
          }
        }
        if !present{
          C.rows = append(C.rows, row3)
        }
      }
    }
  }

  return C
}


//function to merge two item sets if possible
func mergeRow(row1 row, row2 row) (row,bool){
  var row3 row

  var check_item = make([]bool, len(row1.item_set))
  check := 0

  for i,item1 := range row1.item_set{
    check_item[i] = false
    for _,item2 := range row2.item_set{

      if(item1 == item2){
        check++
        check_item[i] = true
      }
    }
  }

  if len(row1.item_set) - check == 1{
    for i,item := range row1.item_set{
      if !check_item[i]{
        row3.item_set = append(row3.item_set, item)
      }
    }

    for _,item := range row2.item_set{
      row3.item_set = append(row3.item_set, item)
    }
    return row3, true
  }

  return row3, false
}

//function to print a table
func printTable(table1 table){
    for _,row1 := range table1.rows{
      for _,item := range row1.item_set{
        fmt.Printf("%2d ", item)
      }
      fmt.Printf("%2d\n", row1.sup_count)
    }

    fmt.Println("\n\n")
}


//function to create L1 from transactions
func createL1(transactions [][]int) table{
  var L1 table

  for _,transaction := range transactions{
    for _,item := range transaction{
      if len(L1.rows) == 0{
        L1.rows = append(L1.rows, row{[]int{item}, 1})
      } else{
        present := false
        for i,row1 := range L1.rows{
          if row1.item_set[0] == item{
            present = true
            L1.rows[i].sup_count++;
          }
        }
        if !present{
          L1.rows = append(L1.rows, row{[]int{item}, 1})
        }
      }
    }
  }

  return L1
}




//function to read the transactions file and load it
func readFile(filename string) [][]int{

  var transactions = make([][]int, 0)

  f, err := os.Open(filename)
  if err != nil{
    fmt.Println(err)
  }
  defer f.Close()

  r := bufio.NewReaderSize(f, 4*1024)

  line, isPrefix, err := r.ReadLine()

  for err == nil && !isPrefix{
    s := string(line)
    slist := strings.Split(s, ",")

    transactions = append(transactions, )

    var transaction []int
    for _,numstring := range slist{
      num, err := strconv.Atoi(numstring)
      if err != nil{
        fmt.Println(err)
      }
      transaction = append(transaction, num)
    }

    transactions = append(transactions, transaction)

    line, isPrefix, err = r.ReadLine()
  }

  return transactions
}
