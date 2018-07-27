
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


  a, aub := getInput()

  sup_count_a := countSupport(a, transactions)
  sup_count_aub := countSupport(aub, transactions)
  fmt.Println("support count of", a, " => ", sup_count_a)
  fmt.Println("support count of", aub, " => ", sup_count_aub)
  fmt.Printf("confidence => %.2f percent\n", (float32(sup_count_aub)/float32(sup_count_a))*100)

}

//function to get input from the user
func getInput() ([]int, []int){

  var a, b []int

  a = make([]int, 0)
  b = make([]int, 0)


  var stringA, stringB  string

  fmt.Println("Enter A:")
  fmt.Scanf("%s ", &stringA)

  fmt.Println("Enter B:", stringB)
  fmt.Scanf("%s ", &stringB)

  alist := strings.Split(stringA, ",")

  for _,numstring := range alist{
      num,err := strconv.Atoi(numstring)
      if err != nil{
        fmt.Println(err)
      }

      a = append(a, num)
  }

  blist := strings.Split(stringB, ",")


    for _,numstring := range blist{
        num,err := strconv.Atoi(numstring)
        if err != nil{
          fmt.Println(err)
        }

        b = append(b, num)
    }

    aub := append(a, b...)

    for i := 0; i < len(aub); i++{
      for j := i+1; j < len(aub);{
        if aub[i] == aub[j]{
        aub[j] = aub[len(aub)-1]
        aub = aub[:len(aub)-1]
        }else{
          j++;
        }
      }
    }

  return a,aub
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
      var row3 row
      var merged bool
      row3.item_set, merged = mergeRow(row1.item_set, row2.item_set)
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
func mergeRow(item_set1 []int , item_set2 []int) ([]int,bool){
  var item_set3 []int

  item_set3 = make([]int, 0)

  var check_item = make([]bool, len(item_set1))
  check := 0

  for i,item1 := range item_set1{
    check_item[i] = false
    for _,item2 := range item_set2{

      if(item1 == item2){
        check++
        check_item[i] = true
      }
    }
  }

  if len(item_set1) - check == 1{
    for i,item := range item_set1{
      if !check_item[i]{
        item_set3 = append(item_set3, item)
      }
    }

    for _,item := range item_set2{
      item_set3 = append(item_set3, item)
    }
    return item_set3, true
  }

  return item_set3, false
}

//function to print a table
func printTable(table1 table){

  if len(table1.rows) == 0{
    return
  }

  for i := 0; i <= len(table1.rows[0].item_set); i++{
    fmt.Printf("----")
  }

  fmt.Println()

  for _,row1 := range table1.rows{
    fmt.Printf("|")
    for _,item := range row1.item_set{
      fmt.Printf("%2d ", item)
    }
    fmt.Printf("|%2d |\n", row1.sup_count)
  }
  for i := 0; i <= len(table1.rows[0].item_set); i++{
    fmt.Printf("----")
  }
  fmt.Println("\n")
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
