package pipeline

import (
	"testing"

	"gopkg.in/check.v1"
)

func TestPipeline(t *testing.T) {
	check.Suite(&PipelineSuite{})
	check.TestingT(t)
}

type PipelineSuite struct {
	pipelineS1 *Pipeline
	pipelineS2 *Pipeline
	pipelineS3 *Pipeline
	pipelineM1 *Pipeline
	pipelineM2 *Pipeline
	pipelineM3 *Pipeline
	pipelineM4 *Pipeline
}

func (s *PipelineSuite) SetUpTest(c *check.C) {
	s.pipelineS1 = NewPipeline()
	s.pipelineS2 = s.pipelineS1.Pipe(stagex("A"))
	s.pipelineS3 = s.pipelineS2.Pipe(stagex("B"))

	s.pipelineM1 = NewPipeline()
	s.pipelineM2 = s.pipelineM1.Pipe(stagex("A1"))
	s.pipelineM3 = s.pipelineM1.Pipe(stagex("A2"))
	s.pipelineM4 = s.pipelineM2.Pipe(stagex("B"))
}

func (s *PipelineSuite) estEmit(c *check.C) {
	rs1 := s.pipelineS1.Listen()
	rs2 := s.pipelineS2.Listen()
	rs3 := s.pipelineS3.Listen()
	rm1 := s.pipelineM1.Listen()
	rm2 := s.pipelineM2.Listen()
	rm3 := s.pipelineM3.Listen()
	rm4 := s.pipelineM4.Listen()

	s.pipelineS1.Emit(Event{Result: ""})
	s.pipelineM1.Emit(Event{Result: ""})

	str1 := (<-rs1).Result.(string)
	str2 := (<-rs2).Result.(string)
	str3 := (<-rs3).Result.(string)
	c.Check(str1, check.Equals, "")
	c.Check(str2, check.Equals, "A")
	c.Check(str3, check.Equals, "AB")

	str1 = (<-rm1).Result.(string)
	str2 = (<-rm2).Result.(string)
	str3 = (<-rm3).Result.(string)
	str4 := (<-rm4).Result.(string)
	c.Check(str1, check.Equals, "")
	c.Check(str2, check.Equals, "A1")
	c.Check(str3, check.Equals, "A2")
	c.Check(str4, check.Equals, "A1B")
}

func (s *PipelineSuite) TestClose(c *check.C) {
	rs1 := s.pipelineS1.Listen()
	rs2 := s.pipelineS2.Listen()
	rs3 := s.pipelineS3.Listen()
	rm1 := s.pipelineM1.Listen()
	rm2 := s.pipelineM2.Listen()
	rm3 := s.pipelineM3.Listen()
	rm4 := s.pipelineM4.Listen()

	s.pipelineS1.Close()
	s.pipelineM1.Close()

	_, ok1 := <-rs1
	_, ok2 := <-rs2
	_, ok3 := <-rs3
	c.Check(ok1, check.Equals, false)
	c.Check(ok2, check.Equals, false)
	c.Check(ok3, check.Equals, false)

	_, ok1 = <-rm1
	_, ok2 = <-rm2
	_, ok3 = <-rm3
	_, ok4 := <-rm4
	c.Check(ok1, check.Equals, false)
	c.Check(ok2, check.Equals, false)
	c.Check(ok3, check.Equals, false)
	c.Check(ok4, check.Equals, false)
}

func stagex(s string) Stage {
	return StageFunc(func(in <-chan Event) <-chan Event {
		out := make(chan Event)

		go func(in <-chan Event, out chan<- Event) {
			for e := range in {
				s = e.Result.(string) + s
				e.Result = s
				out <- e
			}
			close(out)
		}(in, out)

		return out
	})
}
