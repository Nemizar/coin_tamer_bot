package commands

import (
	"context"

	"github.com/Nemizar/coin_tamer_bot/internal/core/domain/models/category"
	"github.com/Nemizar/coin_tamer_bot/internal/core/ports"
	"github.com/Nemizar/coin_tamer_bot/internal/pkg/errs"
)

type defaultCategoryTemplate struct {
	Name     string
	Children []string
}

type CreateDefaultCategoryCommandHandler interface {
	Handle(ctx context.Context, command CreateDefaultCategoryCommand) error
}

var _ CreateDefaultCategoryCommandHandler = createDefaultCategoryCommandHandler{}

type createDefaultCategoryCommandHandler struct {
	logger ports.Logger
	uow    ports.UnitOfWork
}

func NewCreateDefaultCategoryCommandHandler(logger ports.Logger, uow ports.UnitOfWork) (CreateDefaultCategoryCommandHandler, error) {
	if logger == nil {
		return nil, errs.NewValueIsRequiredError("logger")
	}

	if uow == nil {
		return nil, errs.NewValueIsRequiredError("uow")
	}

	return &createDefaultCategoryCommandHandler{
		logger: logger,
		uow:    uow,
	}, nil
}

func (c createDefaultCategoryCommandHandler) Handle(ctx context.Context, command CreateDefaultCategoryCommand) error {
	defer func(uow ports.UnitOfWork) {
		err := uow.RollbackUnlessCommitted()
		if err != nil {
			c.logger.Error("create default category command handler: rollback failed", "err", err)
		}
	}(c.uow)

	err := c.uow.Begin(ctx)
	if err != nil {
		return err
	}

	u, err := c.uow.UserRepository().FindByExternalProvider(ctx, command.Provider(), command.ExternalID())
	if err != nil {
		return err
	}

	for _, tpl := range c.getDefaultsCategory() {
		parent, err := category.New(
			tpl.Name,
			u.ID(),
			nil,
		)
		if err != nil {
			return err
		}

		if err := c.uow.CategoryRepository().Create(ctx, parent); err != nil {
			return err
		}

		for _, childName := range tpl.Children {
			pID := parent.ID()
			child, err := category.New(
				childName,
				u.ID(),
				&pID,
			)
			if err != nil {
				return err
			}

			if err := c.uow.CategoryRepository().Create(ctx, child); err != nil {
				return err
			}
		}
	}

	return c.uow.Commit(ctx)
}

func (c createDefaultCategoryCommandHandler) getDefaultsCategory() []defaultCategoryTemplate {
	return []defaultCategoryTemplate{
		{
			Name: "Покупки",
			Children: []string{
				"🍎 Еда, продукты",
				"👕 Одежда",
				"🏡 Дом, хозяйство",
				"💻 Техника",
			},
		},
		{
			Name: "Обязательные",
			Children: []string{
				"🏠 ЖКХ",
				"📞 Телефон",
				"💸 Налоги",
			},
		},
		{
			Name: "Здоровье",
			Children: []string{
				"🏥 Медицина",
				"🏋️ Спорт, здоровье",
				"💅 Красота",
			},
		},
		{
			Name: "Транспорт",
			Children: []string{
				"🚙 Машина",
				"✈️ Поездки",
			},
		},
		{
			Name: "Прочее",
			Children: []string{
				"🔹 Прочее",
				"⚠️ Внеплановые расходы",
			},
		},
		{
			Name: "Развлечения",
			Children: []string{
				"🎬 Кино, театр",
				"🌍 Путешествия",
				"☕ Кафе",
				"🎁 Сувениры",
			},
		},
		{
			Name: "Праздники",
			Children: []string{
				"🎀 Подарки",
				"🎊 Праздники",
			},
		},
		{
			Name: "Услуги",
			Children: []string{
				"🔧 Услуги/сервисы",
				"🌐 Интернет",
			},
		},
		{
			Name: "Обучение",
			Children: []string{
				"📚 Книги",
				"🎓 Курсы и учеба",
			},
		},
	}
}
