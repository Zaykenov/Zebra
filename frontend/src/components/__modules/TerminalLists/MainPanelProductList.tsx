import { Product } from "@api/terminal";
import ItemCard from "@common/ItemCard";
import { FC } from "react";

interface MainPanelProductListProps {
    mainPanelProducts: Product[],
    checkIfHasModal: (product: Product) => boolean,
    onItemSelect: (product: Product, route: string) => () => void
}

const MainPanelProductList: FC<MainPanelProductListProps> = ({
    mainPanelProducts,
    checkIfHasModal,
    onItemSelect,
}) => {
  return (
    <div className="main-panel-product-list pt-4 flex flex-wrap gap-1 justify-evenly mb-10 pb-4 overflow-x-auto">
      {mainPanelProducts?.map((product) => (
        <div key={product.id} className="main-panel-product-item w-[24%]">
          <ItemCard
            name={product.name}
            cover={product.image}
            price={product.price}
            hasModal={checkIfHasModal(product)}
            onSelect={onItemSelect(product, "mainPanel")}
          />
        </div>
      ))}
    </div>
  );
};

export default MainPanelProductList;
