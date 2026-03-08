func (r *FishRepository) SaveRecipesToCache(fishID string, recipes []fish.Recipe) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, recipe := range recipes {
			// recipe_cacheテーブルへ保存 (Upsert)
			if err := tx.Table("recipe_cache").Save(&recipe).Error; err != nil {
				return err
			}
			
			// fish_recipe_links中間テーブルで紐付け
			link := map[string]interface{}{
				"fishId":   fishID,
				"recipeId": recipe.ID,
				"updatedAt": time.Now(),
			}
			if err := tx.Table("fish_recipe_links").Save(link).Error; err != nil {
				return err
			}
		}
		return nil
	})
}