package main

import (
    "context"
    "database/sql"
    "fmt"
    "log"
    "strconv"

    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/app"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/dialog"
    "fyne.io/fyne/v2/widget"
    "github.com/Vertinska/my-first-app/internal/repository"
    _ "github.com/mattn/go-sqlite3"
)

type fyneApp struct {
    db            *sql.DB
    repo          repository.ProductRepository
    productList   *widget.List
    products      []repository.Product
    modelEntry    *widget.Entry
    companyEntry  *widget.Entry
    priceEntry    *widget.Entry
    statusLabel   *widget.Label
    selectedID    int // Добавляем поле для хранения ID выбранного элемента
}

func main() {
    ctx := context.Background()
    
    // Подключение к БД
    db, err := sql.Open("sqlite3", "store.db")
    if err != nil {
        log.Fatal("Ошибка подключения к БД:", err)
    }
    defer db.Close()

    // Создание таблицы
    createTableSQL := `CREATE TABLE IF NOT EXISTS products (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        model TEXT,
        company TEXT,
        price INTEGER
    );`
    _, err = db.ExecContext(ctx, createTableSQL)
    if err != nil {
        log.Fatal("Ошибка создания таблицы:", err)
    }

    // Создание репозитория
    repo := NewSQLiteRepository(db)

    // Создание Fyne приложения
    a := app.New()
    w := a.NewWindow("Управление товарами")
    w.Resize(fyne.NewSize(600, 500))

    // Инициализация структуры приложения
    appState := &fyneApp{
        db:          db,
        repo:        repo,
        products:    []repository.Product{},
        statusLabel: widget.NewLabel("Готов к работе"),
        selectedID:  -1, // Инициализируем как -1 (ничего не выбрано)
    }

    // Поле ввода модели
    appState.modelEntry = widget.NewEntry()
    appState.modelEntry.SetPlaceHolder("Введите модель...")

    // Поле ввода компании
    appState.companyEntry = widget.NewEntry()
    appState.companyEntry.SetPlaceHolder("Введите компанию...")

    // Поле ввода цены
    appState.priceEntry = widget.NewEntry()
    appState.priceEntry.SetPlaceHolder("Введите цену...")

    // Кнопка добавления товара
    addButton := widget.NewButton("Добавить товар", func() {
        model := appState.modelEntry.Text
        company := appState.companyEntry.Text
        priceStr := appState.priceEntry.Text

        if model == "" || company == "" || priceStr == "" {
            dialog.ShowInformation("Ошибка", "Заполните все поля", w)
            return
        }

        price, err := strconv.Atoi(priceStr)
        if err != nil {
            dialog.ShowInformation("Ошибка", "Цена должна быть числом", w)
            return
        }

        product := repository.Product{
            Model:   model,
            Company: company,
            Price:   price,
        }

        err = repo.SaveProduct(ctx, product)
        if err != nil {
            dialog.ShowInformation("Ошибка", "Не удалось сохранить товар: "+err.Error(), w)
            return
        }

        // Очистка полей
        appState.modelEntry.SetText("")
        appState.companyEntry.SetText("")
        appState.priceEntry.SetText("")

        // Обновление списка
        appState.loadProducts(ctx, w)
        dialog.ShowInformation("Успех", "Товар добавлен", w)
    })

    // Кнопка обновления списка
    refreshButton := widget.NewButton("Обновить список", func() {
        appState.loadProducts(ctx, w)
    })

    // Создание списка товаров
    appState.productList = widget.NewList(
        func() int { return len(appState.products) },
        func() fyne.CanvasObject {
            return widget.NewLabel("Товар")
        },
        func(i widget.ListItemID, o fyne.CanvasObject) {
            label := o.(*widget.Label)
            if i >= 0 && i < len(appState.products) {
                p := appState.products[i]
                label.SetText(fmt.Sprintf("%d. %s (%s) - %d руб.", p.ID, p.Model, p.Company, p.Price))
            }
        },
    )

    // --- ИСПРАВЛЕННАЯ ЧАСТЬ: Обработка выбора в списке ---
    // При выборе элемента сохраняем его ID в структуре appState
    appState.productList.OnSelected = func(id widget.ListItemID) {
        appState.selectedID = int(id)
        fmt.Println("Selected ID:", appState.selectedID) // Для отладки (можно удалить)
    }

    // При снятии выбора сбрасываем ID в -1
    appState.productList.OnUnselected = func(id widget.ListItemID) {
        if appState.selectedID == int(id) {
            appState.selectedID = -1
        }
        fmt.Println("Unselected ID:", id) // Для отладки (можно удалить)
    }

    // Кнопка удаления товара (исправленная версия)
    deleteButton := widget.NewButton("Удалить выбранный", func() {
        if appState.selectedID < 0 || appState.selectedID >= len(appState.products) {
            dialog.ShowInformation("Ошибка", "Сначала выберите товар из списка", w)
            return
        }

        product := appState.products[appState.selectedID]
        
        // TODO: Здесь можно добавить реальное удаление из БД
        // для этого нужно добавить метод DeleteProduct в интерфейс repository.ProductRepository
        // и его реализацию в sqliteRepository
        
        dialog.ShowInformation("Информация", 
            fmt.Sprintf("Выбран товар для удаления:\nID: %d\n%s (%s) - %d руб.", 
                product.ID, product.Model, product.Company, product.Price), w)
    })
    // --- КОНЕЦ ИСПРАВЛЕНИЙ ---

    // Загрузка начальных данных
    appState.loadProducts(ctx, w)

    // Компоновка интерфейса
    form := container.NewVBox(
        widget.NewLabelWithStyle("Добавление нового товара", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
        container.NewGridWithColumns(2,
            widget.NewLabel("Модель:"), appState.modelEntry,
            widget.NewLabel("Компания:"), appState.companyEntry,
            widget.NewLabel("Цена:"), appState.priceEntry,
        ),
        container.NewHBox(addButton, refreshButton),
        widget.NewSeparator(),
        widget.NewLabelWithStyle("Список товаров", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
        appState.productList,
        deleteButton,
        appState.statusLabel,
    )

    w.SetContent(form)
    w.ShowAndRun()
}

func (a *fyneApp) loadProducts(ctx context.Context, w fyne.Window) {
    products, err := a.repo.ListProducts(ctx)
    if err != nil {
        a.statusLabel.SetText("Ошибка загрузки: " + err.Error())
        return
    }
    a.products = products
    a.productList.Refresh()
    a.statusLabel.SetText(fmt.Sprintf("Загружено товаров: %d", len(products)))
}
