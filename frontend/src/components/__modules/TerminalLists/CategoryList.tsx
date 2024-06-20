import { Category } from "@api/terminal";
import ItemCard from "@common/ItemCard";
import { FC } from "react";

interface CategoryListProps {
    categories: any[]
    mainPanelCategoryId: number
    setSelectedCategory: (value: React.SetStateAction<Category | null>) => void
}

const CategoryList: FC<CategoryListProps> = ({
    categories,
    mainPanelCategoryId,
    setSelectedCategory
}) => {
  return (
    <div className="category-item-list flex flex-wrap gap-2 justify-evenly border-b border-gray-400 pb-4 overflow-x-auto">
      {categories?.map(
        (category, idx) =>
          category.id !== mainPanelCategoryId && (
            <div
              key={`${category.id}_${idx}`}
              className="category-item w-[24%]"
            >
              <ItemCard
                name={category.category}
                cover={category.image}
                onSelect={() => setSelectedCategory(category)}
              />
            </div>
          )
      )}
    </div>
  );
};

export default CategoryList