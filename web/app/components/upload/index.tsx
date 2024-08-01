import Image from "next/image";

export default function Sections() {
    return (
        <div className="container mx-auto flex justify-center items-center h-[100vh]">
            <div className="lg:w-[50%] w-full mx-auto text-center">
                <div className="w-full float-start flex justify-center items-center mb-10">
                <Image src="/img/logo.png" width={48} height={48} alt="logo" /> 
                <h1 className="text-[40px] font-semibold">Vultisig</h1>
                </div>
                <div className="w-full float-start bg-[#061B3A] rounded-md h-[500px] p-10">
                    <h2 className="text-2xl font-bold mb-10">Upload your vault share to start</h2>
                    <div className="lg:w-[80%] w-full  p-10 flex justify-center flex-col items-center mx-auto rounded-2xl border border-dashed border-[#33E6BF]">
                        <Image src="/img/qr-code.png" className="mb-10" width={40} height={40} alt="qr" />
                        <h3 className="text-xl font-semibold">Upload your QR code here</h3>
                        <p className="text-sm">
                        Drop your file here or <b className="underline">upload it </b>
                        </p>
                    </div>
                    <p className="lg:w-[80%] mx-auto mt-10">
                    If you didn’t save the QR code yet, you can find it in the app in the top right on the main screen
                    </p>
                    <button disabled className="bg-[#3a51743b] mt-5 h-[44px] w-[80%] mx-auto rounded-3xl">
                        Start  
                    </button>
                </div>
                <div className="w-full float-start mt-10 text-center">
                    <p>Don’t have a vault yet? Download Vultisig now</p>
                </div>
                <div className="w-full gap-5 float-start flex mt-10 justify-center items-center">
                    <div>
                    <Image src="/img/apple.png" className="mb-10" width={125} height={48} alt="iPhone" />
                    <p>iPhone</p>
                    </div>
                    <div>
                    <Image src="/img/google.png" className="mb-10" width={125} height={48} alt="Android" />
                    <p>Android</p>
                    </div>
                    <div>
                    <Image src="/img/git.png" className="mb-10 h-[40px]" width={125} height={30} alt="Mac" />
                    <p>Mac</p>
                    </div>
                </div>
            </div>
        </div>
    );
}