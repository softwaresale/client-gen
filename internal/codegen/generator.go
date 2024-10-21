package codegen

type ITargetFormatter interface {
	Format(service CompiledService) error
}
