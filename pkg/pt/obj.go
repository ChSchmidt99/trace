package pt

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

//TODO: Move material somewhere else, or parse .mat file
func ParseFromPath(path string, material Material) *Mesh {
	objFile, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer objFile.Close()
	faces := parseOBJ(objFile, material)
	return newMesh(faces)
}

func parseOBJ(objFile *os.File, material Material) []Primitive {
	scanner := bufio.NewScanner(objFile)

	vertecies := make([]Vector3, 1, 1024)
	normals := make([]Vector3, 1, 1024)
	triangles := make([]Primitive, 0, 1024)

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		key := fields[0]
		values := fields[1:]

		switch key {
		case "v":
			if numbers, err := parseFloat(values); err != nil {
				panic(err)
			} else {
				vertecies = append(vertecies, NewVector3(numbers[0], numbers[1], numbers[2]))
			}
		case "vt":
			// TODO: Implement when needed
		case "vn":
			if numbers, err := parseFloat(values); err != nil {
				panic(err)
			} else {
				normals = append(normals, NewVector3(numbers[0], numbers[1], numbers[2]))
			}
		case "f":
			if face, err := parseFace(values, vertecies, normals, material); err != nil {
				panic(err)
			} else {
				triangles = append(triangles, face...)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	//fmt.Printf("Tris: %v\n", len(triangles))

	return triangles
}

func parseFace(args []string, vertecies []Vector3, normals []Vector3, material Material) ([]Primitive, error) {
	vIndeces := make([]int, 0, 4)
	nIndeces := make([]int, 0, 4)
	for _, arg := range args {
		indeces := strings.Split(arg, "/")
		vIndex, err := strconv.Atoi(indeces[0])
		if err != nil {
			return nil, err
		}
		if vIndex < 0 {
			vIndex = len(vertecies) + vIndex
		}
		vIndeces = append(vIndeces, vIndex)

		if len(indeces) >= 3 {
			nIndex, err := strconv.Atoi(indeces[2])
			if err != nil {
				return nil, err
			}
			if nIndex < 0 {
				nIndex = len(normals) + nIndex
			}
			nIndeces = append(nIndeces, nIndex)
		}
	}

	if len(nIndeces) == len(args) {
		triangles := make([]Primitive, 0)
		for i := 1; i+2 <= len(vIndeces); i++ {
			triangles = append(triangles, triangleForIndeces(append(vIndeces[0:1], vIndeces[i:i+2]...), append(nIndeces[0:1], nIndeces[i:i+2]...), vertecies, normals, material))
		}
		return triangles, nil
	} else {
		triangles := make([]Primitive, 0)
		for i := 1; i+2 <= len(vIndeces); i++ {
			triangles = append(triangles, triangleWithoutNormals(append(vIndeces[0:1], vIndeces[i:i+2]...), vertecies, material))
		}
		return triangles, nil
	}
}

func triangleWithoutNormals(vIndeces []int, vertecies []Vector3, material Material) *triangle {
	return newTriangleWithoutNormals(vertecies[vIndeces[0]], vertecies[vIndeces[1]], vertecies[vIndeces[2]], material)
}

func triangleForIndeces(vIndeces []int, nIndeces []int, vertecies []Vector3, normals []Vector3, material Material) *triangle {
	var v [3]vertex
	v[0] = vertex{
		position: vertecies[vIndeces[0]],
		normal:   normals[0],
	}
	v[1] = vertex{
		position: vertecies[vIndeces[1]],
		normal:   normals[1],
	}
	v[2] = vertex{
		position: vertecies[vIndeces[2]],
		normal:   normals[2],
	}
	return newTriangle(v, material)
}

func parseFloat(args []string) ([]float64, error) {
	result := make([]float64, 0, len(args))
	for _, arg := range args {
		if num, err := strconv.ParseFloat(arg, 64); err != nil {
			return nil, err
		} else {
			result = append(result, num)
		}
	}
	return result, nil
}
