package ada

import (
	"fmt"
	"testing"
)

func TestEstimate(t *testing.T) {
	/*
		RangeQueries rangeQ = new RangeQueries(11, 10, 1.004, SketchType.CMS, "exp", 4);
		rangeQ.update(1,1,5);
		rangeQ.update(1, 2, 7);
		rangeQ.update(1, 3, 11);
		rangeQ.update(1, 6, 12);
		rangeQ.update(1, 4, 100);
		rangeQ.update(1, 5, 11);
		rangeQ.update(1, 7, 1);
		rangeQ.update(1, 8, 5);
		rangeQ.update(1, 10, 5);


		Console.WriteLine(rangeQ.estimateinRange(1, 4, 8));
		Console.WriteLine(rangeQ.estimateinRange(1, 1, 10));
		Console.WriteLine(rangeQ.estimateinRange(1, 4, 9));
		Console.WriteLine(rangeQ.estimateinRange(1, 2, 7));
		Console.WriteLine(rangeQ.estimateinRange(1, 2, 5));
		Console.WriteLine(rangeQ.estimateinRange(1, 2, 2));
		Console.WriteLine(rangeQ.estimateinRange(1, 4, 4));
		Console.ReadKey();
	*/
	sks := NewSketches(11, 10, 4, 1.004)
	item1 := []byte("foo")
	sks.Update(item1, 1, 5)
	sks.Update(item1, 2, 7)
	sks.Update(item1, 3, 11)
	sks.Update(item1, 6, 12)
	sks.Update(item1, 4, 100)
	sks.Update(item1, 5, 11)
	sks.Update(item1, 7, 1)
	sks.Update(item1, 8, 5)
	sks.Update(item1, 10, 5)
	fmt.Println(sks.Estimate(item1, 4, 8))

	item2 := []byte("bar")
	sks.Update(item2, 11, 500)
	sks.Update(item2, 5, 11)
	sks.Update(item2, 7, 1)
	sks.Update(item2, 8, 5)
	sks.Update(item2, 10, 510)

	fmt.Println(sks.Estimate(item1, 4, 8))
	fmt.Println(sks.Estimate(item2, 10, 11))
}
