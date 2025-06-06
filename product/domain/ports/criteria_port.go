package ports

type CriterionPort interface {
	And(...CriterionPort) CriterionPort
	Or(...CriterionPort) CriterionPort
}

type CriterionBuilderPort interface {
	Where(entity string, property string, operator string, value any) CriterionPort
	And(...CriterionPort) CriterionPort
	Or(...CriterionPort) CriterionPort
}
