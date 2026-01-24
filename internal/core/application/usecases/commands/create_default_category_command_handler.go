package commands

import (
	"context"

	"github.com/Nemizar/coin_tamer_bot/internal/core/domain/models/category"
	"github.com/Nemizar/coin_tamer_bot/internal/core/ports"
	"github.com/Nemizar/coin_tamer_bot/internal/pkg/errs"
)

type defaultCategoryTemplate struct {
	name     string
	children []string
	cType    category.Type
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

	hasCategories, err := c.uow.CategoryRepository().HasCategoriesByUserID(ctx, u.ID())
	if err != nil {
		return err
	}

	if hasCategories {
		return errs.NewEntityAlreadyExistsError("categories", "user_id", u.ID().String())
	}

	for _, tpl := range c.getDefaultsCategory() {
		parent, err := category.New(
			tpl.name,
			tpl.cType,
			u.ID(),
			nil,
		)
		if err != nil {
			return err
		}

		if err := c.uow.CategoryRepository().Create(ctx, parent); err != nil {
			return err
		}

		for _, childName := range tpl.children {
			pID := parent.ID()
			child, err := category.New(
				childName,
				tpl.cType,
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
			name: "Покупки",
			children: []string{
				"🍎 Еда, продукты",
				"👕 Одежда",
				"🏡 Дом, хозяйство",
				"💻 Техника",
			},
			cType: category.TypeExpense,
		},
		{
			name: "Обязательные",
			children: []string{
				"🏠 ЖКХ",
				"📞 Телефон",
				"💸 Налоги",
			},
			cType: category.TypeExpense,
		},
		{
			name: "Здоровье",
			children: []string{
				"🏥 Медицина",
				"🏋️ Спорт, здоровье",
				"💅 Красота",
			},
			cType: category.TypeExpense,
		},
		{
			name: "Транспорт",
			children: []string{
				"🚙 Машина",
				"✈️ Поездки",
			},
			cType: category.TypeExpense,
		},
		{
			name: "Прочее",
			children: []string{
				"🔹 Прочее",
				"⚠️ Внеплановые расходы",
			},
			cType: category.TypeExpense,
		},
		{
			name: "Развлечения",
			children: []string{
				"🎬 Кино, театр",
				"🌍 Путешествия",
				"☕ Кафе",
				"🎁 Сувениры",
			},
			cType: category.TypeExpense,
		},
		{
			name: "Праздники",
			children: []string{
				"🎀 Подарки",
				"🎊 Праздники",
			},
			cType: category.TypeExpense,
		},
		{
			name: "Услуги",
			children: []string{
				"🔧 Услуги/сервисы",
				"🌐 Интернет",
			},
			cType: category.TypeExpense,
		},
		{
			name: "Обучение",
			children: []string{
				"📚 Книги",
				"🎓 Курсы и учеба",
			},
			cType: category.TypeExpense,
		},
		{
			name:  "Зарплата",
			cType: category.TypeIncome,
		},
		{
			name:  "Проценты с вклада",
			cType: category.TypeIncome,
		},
		{
			name:  "Кешбек",
			cType: category.TypeIncome,
		},
	}
}
