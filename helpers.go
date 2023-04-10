package grawler

import "github.com/DAtek/gotils"

func joinPath(parts ...string) string {
	return gotils.Reduce(parts, joinTwoPaths)
}

func joinTwoPaths(path1, path2 string) string {
	lenPath1 := len(path1)

	if string(path1[lenPath1-1]) == "/" {
		path1 = path1[:lenPath1-1]
	}

	if string(path2[0]) == "/" {
		path2 = path2[1:]
	}

	return path1 + "/" + path2

}
