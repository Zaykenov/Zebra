import { Carousel, IconButton } from "@material-tailwind/react";
import {
  ArrowLeftIcon,
  ArrowLeftOnRectangleIcon,
  ArrowRightIcon,
} from "@heroicons/react/24/outline";
import { FC, useState } from "react";
import Image from "next/image";
export interface InstructionImage {
  slideImagePath: string;
  width: number;
  height: number;
}

export interface InstructionSlide {
  introText?: string;
  headerComponentElement?: JSX.Element;
  slideImage?: InstructionImage;
  footerComponentElement?: JSX.Element;
}

interface InstructionCarouselProps {
  slides: InstructionSlide[];
}

const InstructionCarousel: FC<InstructionCarouselProps> = ({ slides }) => {
  const [currentSlide, setCurrentSlide] = useState<number>(0)
  return (
    // <Carousel
    //   className="rounded-xl"
    //   prevArrow={({ handlePrev }) => (
    //     <IconButton
    //       variant="text"
    //       size="lg"
    //       onClick={handlePrev}
    //       className="!absolute bottom-2 left-4"
    //     >
    //       <ArrowLeftIcon strokeWidth={2} className="w-6 h-6"/>
    //     </IconButton>
    //   )}
    //   nextArrow={({ handleNext }) => {
    //     return (
    //       <IconButton
    //         variant="text"
    //         size="lg"
    //         onClick={handleNext}
    //         className="!absolute bottom-2 right-4"
    //       >
    //         <ArrowRightIcon strokeWidth={2} className="w-6 h-6"/>
    //       </IconButton>
    //     );
    //   }}
    // >
    //   {slides.map((slide, idx) => {
    //     if (slide.introText) {
    //       return (
    //         <div key={idx} className="flex-col align-center mt-10">
    //           <div className="flex justify-center gap-4">
    //             <h1 className="text-8xl font-bold text-primary mr-2">
    //               {slide.introText}
    //             </h1>
    //           </div>
    //         </div>
    //       );
    //     } else
    //       return (
    //         <div key={idx} className="flex-col align-center mt-10">
    //           <div className="flex justify-center gap-4">
    //             <h1 className="text-8xl font-bold text-primary mr-2">{idx % 5}</h1>
    //             {slide.headerComponentElement}
    //           </div>
    //           <div className="mt-9 flex justify-center">
    //             <Image
    //               src={`/images/${slide.slideImage?.slideImagePath}`}
    //               width={slide.slideImage?.width}
    //               height={slide.slideImage?.height}
    //               alt="InstructionImage"
    //             />
    //           </div>
    //           <div className="mt-9 flex justify-center">
    //             <p className="text-xl font-semibold">
    //               {slide.footerComponentElement ? (
    //                 slide.footerComponentElement
    //               ) : (
    //                 <button 
    //                   className="text-white shadow-md pt-1.5 pb-2 px-3 bg-primary text-sm font-semibold rounded-md hover:bg-teal-600"
    //                   onClick={()=>{

    //                   }}
    //                 >
    //                   Посмотреть заново
    //                 </button>
    //               )}
    //             </p>
    //           </div>
    //         </div>
    //       );
    //   })}
    // </Carousel>
    <>
    
    </>
  );
};

export default InstructionCarousel;
