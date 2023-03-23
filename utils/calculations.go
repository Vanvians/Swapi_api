package utils

import "fmt"


func ConvertCmToFtIn(totalHeightCm int) (string, error) {
    // Convert the specified height in centimeters to feet and inches
    totalHeightInches := float64(totalHeightCm) / 2.54
    feet := int(totalHeightInches / 12)
    inches := int(totalHeightInches) % 12
    
    return fmt.Sprintf("%d ft %d in", feet, inches), nil
}