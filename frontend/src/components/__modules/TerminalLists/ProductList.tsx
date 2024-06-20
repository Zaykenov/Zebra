import { Category, Product } from "@api/terminal"
import ItemCard from "@common/ItemCard"
import { FC } from "react"

interface ProductListProps {
    selectedCategory: Category,
    checkIfHasModal: (product: Product) => boolean,
    onItemSelect: (product: Product, route: string) => () => void
}

const ProductList: FC<ProductListProps> = ({
    selectedCategory,
    checkIfHasModal,
    onItemSelect
}) => {
    return (
        <div className="product-item-list flex flex-wrap gap-2 justify-evenly pb-4 overflow-x-auto">
        {selectedCategory.products?.map((product) => (
          <div key={product.id} className="product-item w-[24%]">
            <ItemCard
              name={`${product.name}`}
              cover={product.image}
              price={product.price}
              hasModal={checkIfHasModal(product)}
              onSelect={onItemSelect(product, "categories")}
            />
          </div>
        ))}
      </div>
    )
}

export default ProductList